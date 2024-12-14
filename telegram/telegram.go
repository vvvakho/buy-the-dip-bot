package telegram

import (
	"buy-the-dip-bot/internal/db"
	"buy-the-dip-bot/utils"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

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

func ListenForUpdates(queriesDB *db.Queries) {
	if bot == nil {
		log.Fatal("Bot is not initialized.")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleCommand(update.Message, queriesDB)
		}
	}
}

func SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
