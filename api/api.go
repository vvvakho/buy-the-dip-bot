package api

import (
	"buy-the-dip-bot/telegram"
	"fmt"
	"log"
)

type MarketDataProvider interface {
	FetchRSI(ticker string) (float64, error)
}

func TrackRSI(ticker string, m MarketDataProvider) {
	rsi, err := m.FetchRSI(ticker)
	if err != nil {
		log.Fatalf("Unable to track RSI: %v", err)
	}

	tgSendRSI(ticker, rsi)
}

func tgSendRSI(ticker string, rsi float64) error {
	// tgChatID will later be replaced by a slice of subscribed users' IDs
	tgChatID := telegram.ChatID
	telegram.SendMessage(tgChatID, fmt.Sprintf("RSI of %s: %f", ticker, rsi))
	return nil
}
