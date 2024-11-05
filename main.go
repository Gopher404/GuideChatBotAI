package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/config"
	"main/database"
	"main/gigachat"
	"main/tg_bot"
	"main/usecase"
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

	db, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	attractionsRepo := database.NewAttractionsRepo(db)

	AIClient, err := gigachat.NewClient(&cfg.GigaChat)
	if err != nil {
		log.Fatal(err)
	}

	tourMaker := usecase.NewTourMaker(attractionsRepo, AIClient)

	handler := tg_bot.NewHandler(tourMaker)

	bot, err := tg_bot.NewBot(handler, &cfg.TGBot)
	if err != nil {
		log.Fatal(err)
	}
	bot.Run()
}
