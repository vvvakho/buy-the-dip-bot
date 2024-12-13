package telegram

import (
	"buy-the-dip-bot/utils"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

// temporarily storing chatID here before user storage logic
var ChatID int64

func InitBot() error {
	var err error

	token, err := utils.GetEnv("TELEGRAM_BOT_TOKEN")
	if err != nil {
		log.Fatal("Error loading env variable")
	}
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
		SendMessage(message.Chat.ID, "Welcome to the bot!")
	case "/subscribe":
		SendMessage(message.Chat.ID, "You've subscribed!")
	case "/unsubscribe":
		SendMessage(message.Chat.ID, "You've unsubscribed!")
	default:
		SendMessage(message.Chat.ID, "Unknown command.")
	}
}

func SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
