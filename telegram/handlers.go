package telegram

import (
	"buy-the-dip-bot/internal/db"
	"context"
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func handleCommand(message *tgbotapi.Message, queriesDB *db.Queries) {
	switch message.Text {
	case "/start":
		SendMessage(message.Chat.ID, "Welcome to the bot!")
	case "/subscribe":
		// TODO: change to customizable ticker
		ticker := "SPY"
		if err := AddSubscriptionDB(message.Chat.ID, ticker, queriesDB); err != nil {
			log.Printf("Error adding subscription: %v", err)
			SendMessage(message.Chat.ID, "Something went wrong. Please try again later.")
			return
		}
		SendMessage(message.Chat.ID, "You've subscribed!")
	case "/unsubscribe":
		SendMessage(message.Chat.ID, "You've unsubscribed!")
	default:
		SendMessage(message.Chat.ID, "Unknown command.")
	}
}

func AddSubscriptionDB(userID int64, ticker string, queriesDB *db.Queries) error {
	getSubParams := db.GetSubscriptionParams{
		UserID: userID,
		Ticker: ticker,
	}
	subRow, err := queriesDB.GetSubscription(context.Background(), getSubParams)
	if err != nil {
		if err == sql.ErrNoRows {
			// no subscription present
			addSubParams := db.AddSubscriptionParams{
				SubID:  uuid.New(),
				UserID: userID,
				Ticker: ticker,
			}
			if err := queriesDB.AddSubscription(context.Background(), addSubParams); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if subRow.Active.Valid && subRow.Active.Bool {
		return nil
	} else {
		reacSubParams := db.ReactivateSubscriptionParams{
			UserID: userID,
			Ticker: ticker,
		}
		if err := queriesDB.ReactivateSubscription(context.Background(), reacSubParams); err != nil {
			return err
		}
	}

	return nil
}
