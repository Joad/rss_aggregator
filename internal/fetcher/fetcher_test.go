package fetcher

import (
	"testing"
)

func TestParseXml(t *testing.T) {
	data := `<rss>
		<channel>
			<title>Test</title>
			<item>
				<title>Test Item</title>
			</item>
		</channel>
	</rss>`
	result, err := parseXml([]byte(data))
	if err != nil {
		t.Fatal("Parsing failed: ", err)
	}
	if result.Channel.Title != "Test" {
		t.Fatal("Incorrect title: ", result.Channel.Title)
	}
	itemTitle := result.Channel.Items[0].Title
	if itemTitle != "Test Item" {
		t.Fatal("Incorrect item title: ", itemTitle)
	}
}

func TestParseMultipleItems(t *testing.T) {
	data := `<rss>
		<channel>
			<title>Test</title>
			<item>
				<title>Test Item 1</title>
			</item>
			<item>
				<title>Test Item 2</title>
			</item>
		</channel>
	</rss>`
	result, err := parseXml([]byte(data))
	if err != nil {
		t.Fatal("Parsing failed: ", err)
	}
	if result.Channel.Title != "Test" {
		t.Fatal("Incorrect title: ", result.Channel.Title)
	}
	if len(result.Channel.Items) != 2 {
		t.Fatal("Incorrect number of items")
	}
	itemTitle := result.Channel.Items[0].Title
	if itemTitle != "Test Item 1" {
		t.Fatal("Incorrect item title: ", itemTitle)
	}
}
