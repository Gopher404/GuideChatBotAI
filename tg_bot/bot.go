package tg_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/config"
	"strings"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	MessageHandler func(update *tgbotapi.Update) (tgbotapi.Chattable, error)
	cfg            *config.TGBotConfig
}

func NewBot(cfg *config.TGBotConfig) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	api.Debug = cfg.Debug

	log.Printf("Authorized on account %s", api.Self.UserName)

	bot := &Bot{
		api: api,
		cfg: cfg,
	}
	bot.MessageHandler = MessageHandler

	return bot, nil
}

func (bot *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = bot.cfg.Timeout

	updates := bot.api.GetUpdatesChan(u)

	var tokens chan struct{}

	if bot.cfg.MaxThreads >= 0 {
		tokens = make(chan struct{}, bot.cfg.MaxThreads) // ограничение количества горутин
	}

	for update := range updates {
		if update.Message != nil {
			go func() {
				if bot.cfg.MaxThreads >= 0 {
					tokens <- struct{}{}
					defer func() { <-tokens }()
				}

				msg, err := bot.MessageHandler(&update)
				if err != nil {
					log.Printf("error handle message (%s)\n: %s", strings.ReplaceAll(update.Message.Text, "\n", "\\n"), err.Error())
				}

				if _, err := bot.api.Send(msg); err != nil {
					log.Println(err)
				}
			}()
		}
	}
}
