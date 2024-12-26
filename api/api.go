package api

import (
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/telegram"
	"buy-the-dip-bot/utils"
	"context"
	"fmt"
	"log"
	"time"
)

type MarketDataProvider interface {
	FetchRSI(ticker string, date time.Time, queriesDB *db.Queries) (RSI, error)
	FetchDaily(ticker string, date time.Time, queriesDB *db.Queries) (Daily, error)
}

func TrackRSI(ticker string, m MarketDataProvider, queriesDB *db.Queries) {
	rsi, err := m.FetchRSI(ticker, time.Now(), queriesDB)
	if err != nil {
		log.Printf("Unable to track RSI of %s: %v", ticker, err)
	}

	tgSendRSI(ticker, RSI{RSI: rsi.RSI, Date: rsi.Date}, queriesDB)
}

func TrackPrice(ticker string, m MarketDataProvider, queriesDB *db.Queries) {
	daily, err := m.FetchDaily(ticker, time.Now(), queriesDB)
	if err != nil {
		log.Printf("Unable to track price of %s: %v", ticker, err)
	}

	tgSendDaily(ticker, daily, time.Now(), queriesDB)
}

func tgSendRSI(ticker string, rsi RSI, queriesDB *db.Queries) error {
	subRows, err := queriesDB.GetSubscriptionsByTicker(context.Background(), ticker)
	if err != nil {
		return err
	}

	for _, sub := range subRows {
		comment := utils.CommentRSI(rsi.RSI)
		formattedMessage := utils.FormatMessage(rsi.Date, ticker, rsi.RSI, comment)
		telegram.SendMessage(sub.UserID, formattedMessage)
	}

	return nil
}

func tgSendDaily(ticker string, daily Daily, date time.Time, queriesDB *db.Queries) {
	notifyText := fmt.Sprintf(`📈 Daily Stock Update: %s (as of %v)
	- Close Price: $%.2f
	- Volume: %d
	- Price Range: $%.2f - $%.2f`,
		ticker,
		date,
		daily.Close,
		daily.Volume,
		daily.Low,
		daily.High,
	)

	telegram.NotifyAllUsers(queriesDB, notifyText)
}
