package tg_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TourMaker interface {
	NewTour(userMessage string) (string, error)
}

type Handler struct {
	tourMaker TourMaker
}

func NewHandler(tourMaker TourMaker) *Handler {
	return &Handler{
		tourMaker: tourMaker,
	}
}

func (h *Handler) MessageHandler(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	tour, err := h.tourMaker.NewTour(update.Message.Text)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, tour)

	return msg, nil
}
