package telegram

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var bot *tgbotapi.BotAPI

func InitBot() error {
	var err error

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Token not found")
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return nil
}

func ListenForUpdates() {
	if bot == nil {
		log.Fatal("Bot is not initialized.")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleCommand(update.Message)
		}
	}
}

func handleCommand(message *tgbotapi.Message) {
	switch message.Text {
	case "/start":
		sendMessage(message.Chat.ID, "Welcome to the bot!")
	case "/subscribe":
		sendMessage(message.Chat.ID, "You've subscribed!")
	case "/unsubscribe":
		sendMessage(message.Chat.ID, "You've unsubscribed!")
	default:
		sendMessage(message.Chat.ID, "Unknown command.")
	}
}

func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
