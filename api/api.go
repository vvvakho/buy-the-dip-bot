package api

import (
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/telegram"
	"buy-the-dip-bot/utils"
	"context"
	"log"
	"time"
)

type MarketDataProvider interface {
	FetchRSI(ticker string, date time.Time, queriesDB *db.Queries) (RSI, error)
}

func TrackRSI(ticker string, m MarketDataProvider, queriesDB *db.Queries) {
	rsi, err := m.FetchRSI(ticker, time.Now(), queriesDB)
	if err != nil {
		log.Fatalf("Unable to track RSI: %v", err)
	}

	tgSendRSI(ticker, RSI{RSI: rsi.RSI, Date: rsi.Date}, queriesDB)
}

func tgSendRSI(ticker string, rsi RSI, queriesDB *db.Queries) error {
	subRows, err := queriesDB.GetSubscriptionsByTicker(context.Background(), ticker)
	if err != nil {
		return err
	}

	for _, sub := range subRows {
		formattedMessage := utils.FormatMessage(rsi.Date, ticker, rsi.RSI)
		telegram.SendMessage(sub.UserID, formattedMessage)
	}

	return nil
}
