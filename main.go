package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/config"
	"main/tg_bot"
	"os"
)

const configPath = "config/config.json"

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	if err := tgbotapi.SetLogger(log.New(os.Stdout, "Bot: ", log.Ltime)); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Read(configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	bot, err := tg_bot.NewBot(&cfg.TGBot)
	if err != nil {
		log.Fatal(err)
	}
	bot.Run()

}
