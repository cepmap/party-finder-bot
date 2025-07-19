package main

import (
	"log"
	"time"

	"github.com/cepmap/party-finder-bot/internal/config"
	"github.com/cepmap/party-finder-bot/internal/dbconnector"
	"github.com/cepmap/party-finder-bot/internal/handlers"
	tele "gopkg.in/telebot.v4"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load(".env")

	db, err := dbconnector.NewDBConnectorFromConfig(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Проверяем и создаем админов в БД
	if err := db.EnsureAdminsExists(cfg.AdminTelegramIDs); err != nil {
		log.Printf("Warning: failed to ensure admins exist: %v", err)
	}

	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	h := handlers.NewHandlers(db, cfg)
	h.Register(b)

	log.Printf("Bot started")
	b.Start()
}
