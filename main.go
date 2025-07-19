package main

import (
	"log"
	"os"
	"time"

	"github.com/cepmap/party-finder-bot/internal/handlers"
	tele "gopkg.in/telebot.v4"
)

func main() {
	log.Println("Starting bot")
	token, err := os.ReadFile("token")
	if err != nil {
		log.Fatal(err)
		return
	}
	pref := tele.Settings{
		Token:  string(token),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	h := handlers.NewHandlers()
	h.Register(b)

	b.Start()
}
