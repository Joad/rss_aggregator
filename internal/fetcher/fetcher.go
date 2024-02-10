package fetcher

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Joad/rss_aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FeedData struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeeds(db *database.Queries) {
	log.Println("Fetcher started")
	for {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), 3)
		if err != nil {
			log.Println("Error fetching feeds: ", err)
		}
		log.Printf("Going to fetch %d feeds\n", len(feeds))
		wg := &sync.WaitGroup{}
		wg.Add(len(feeds))
		for _, feed := range feeds {
			go scrapeFeed(db, wg, feed)
		}
		log.Println("Waiting for feeds to finish")
		wg.Wait()
		log.Println("All feeds done")
		time.Sleep(time.Minute)
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched: %v\n", feed.Name, err)
		return
	}
	feedData, err := fetchFeedData(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
		return
	}
	for _, item := range feedData.Channel.Items {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				if err.Code == "23505" {
					continue
				}
			}
			log.Println("Error saving post: ", err)
		}
	}
	log.Printf("Fetched %s, item count: %d\n", feedData.Channel.Title, len(feedData.Channel.Items))
}

func fetchFeedData(url string) (FeedData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return FeedData{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return FeedData{}, err
	}
	return parseXml(data)
}

func parseXml(data []byte) (FeedData, error) {
	result := FeedData{}
	err := xml.Unmarshal(data, &result)
	return result, err
}
