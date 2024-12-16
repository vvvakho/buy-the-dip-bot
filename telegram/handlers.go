package telegram

import (
	"buy-the-dip-bot/database"
	"buy-the-dip-bot/internal/db"
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ticker string = "SPY"
var userErrorMsg string = "Something went wrong. Please try again later."

func handleCommand(message *tgbotapi.Message, queriesDB *db.Queries) {
	switch message.Text {
	case "/start":
		// TODO : add user to usersDB
		SendMessage(message.Chat.ID, "Welcome to the bot!")
	case "/subscribe":
		handleSubscribe(message, queriesDB)
	case "/unsubscribe":
		handleUnsubscribe(message, queriesDB)
	default:
		SendMessage(message.Chat.ID, "Unknown command.")
	}
}

func handleSubscribe(message *tgbotapi.Message, queriesDB *db.Queries) {
	// TODO: change to customizable ticker
	if err := database.AddSubscription(message.Chat.ID, ticker, queriesDB); err != nil {
		if errors.Is(err, database.ErrAlreadySubscribed) {
			SendMessage(message.Chat.ID, "Already subscribed!")
			return
		}

		log.Printf("Error adding subscription: %v", err)
		SendMessage(message.Chat.ID, userErrorMsg)
		return
	}
	SendMessage(message.Chat.ID, "You've subscribed!")
}

func handleUnsubscribe(message *tgbotapi.Message, queriesDB *db.Queries) {
	if err := database.Unsubscribe(message.Chat.ID, ticker, queriesDB); err != nil {
		log.Printf("Error unsubscribing: %v", err)
		SendMessage(message.Chat.ID, userErrorMsg)
		return
	}
	SendMessage(message.Chat.ID, "You've unsubscribed!")
}
