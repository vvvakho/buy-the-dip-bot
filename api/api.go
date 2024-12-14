package api

import (
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/telegram"
	"context"
	"fmt"
	"log"
)

type MarketDataProvider interface {
	FetchRSI(ticker string) (float64, error)
}

func TrackRSI(ticker string, m MarketDataProvider, queriesDB *db.Queries) {
	rsi, err := m.FetchRSI(ticker)
	if err != nil {
		log.Fatalf("Unable to track RSI: %v", err)
	}

	tgSendRSI(ticker, rsi, queriesDB)
}

func tgSendRSI(ticker string, rsi float64, queriesDB *db.Queries) error {
	subRows, err := queriesDB.GetSubscriptionsByTicker(context.Background(), ticker)
	if err != nil {
		return err
	}

	for _, sub := range subRows {
		telegram.SendMessage(sub.UserID, fmt.Sprintf("RSI of %s: %f", ticker, rsi))
	}

	return nil
}
