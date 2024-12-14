package main

import (
	"buy-the-dip-bot/api"
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/telegram"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// telegram client
	if err := telegram.InitBot(); err != nil {
		log.Fatalf("Unable to initialize telegram bot: %v", err)
	}

	// database
	godotenv.Load()
	dbURL := os.Getenv("POSTGRES_URL")
	postgres, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to initialize database: %v", err)
	}
	queriesDB := db.New(postgres)

	// alpha vantage client
	go func() {
		av, err := api.InitAlphaVantageClient()
		if err != nil {
			log.Printf("Unable to initialize alpha vantage client: %v", err)
		}

		for {
			api.TrackRSI("SPY", av, queriesDB)
			time.Sleep(15 * time.Second)
		}

	}()

	// routine for listening for telegram
	for {
		telegram.ListenForUpdates(queriesDB)
	}

}
