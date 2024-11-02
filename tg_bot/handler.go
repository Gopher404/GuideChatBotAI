package tg_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func MessageHandler(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "-")

	return msg, nil
}
