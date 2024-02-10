package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Joad/rss_aggregator/internal/database"
	"github.com/Joad/rss_aggregator/internal/fetcher"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}

	go fetcher.FetchFeeds(cfg.DB)

	server := http.Server{
		Handler: cfg.routes(),
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
