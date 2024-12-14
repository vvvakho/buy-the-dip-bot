package main

import (
	"buy-the-dip-bot/api"
	"buy-the-dip-bot/telegram"
	"log"
	"time"
)

func main() {
	// telegram client
	if err := telegram.InitBot(); err != nil {
		log.Fatalf("Unable to initialize telegram bot: %v", err)
	}

	// database
	//
	// goroutine for api logic
	//
	//

	// alpha vantage client
	go func() {
		av, err := api.InitAlphaVantageClient()
		if err != nil {
			log.Printf("Unable to initialize alpha vantage client: %v", err)
		}

		for {
			api.TrackRSI("SPY", av)
			time.Sleep(15 * time.Second)
		}

	}()

	// routine for listening for telegram
	for {
		telegram.ListenForUpdates()
	}

}
