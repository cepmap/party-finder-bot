package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Handlers struct {
	MainMarkup *tele.ReplyMarkup
	SignupsBtn tele.Btn
	GamesBtn   tele.Btn
	EventsBtn  tele.Btn
}

func NewHandlers() *Handlers {
	mainMarkup := &tele.ReplyMarkup{ResizeKeyboard: true}
	signupsBtn := mainMarkup.Text("Мои записи на события")
	gamesBtn := mainMarkup.Text("Мои игры")
	eventsBtn := mainMarkup.Text("Мои события")
	mainMarkup.Reply(
		mainMarkup.Row(signupsBtn),
		mainMarkup.Row(gamesBtn, eventsBtn),
	)

	return &Handlers{
		MainMarkup: mainMarkup,
		SignupsBtn: signupsBtn,
		GamesBtn:   gamesBtn,
		EventsBtn:  eventsBtn,
	}

}

func (h *Handlers) Register(b *tele.Bot) {
	b.Handle("/start", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("Привет, %s", c.Sender().Username), h.MainMarkup)
	})
	b.Handle(&h.SignupsBtn, func(c tele.Context) error {
		return c.Send("1", h.MainMarkup)
	})
	b.Handle(&h.GamesBtn, func(c tele.Context) error {
		return c.Send("2", h.MainMarkup)
	})
	b.Handle(&h.EventsBtn, func(c tele.Context) error {
		return c.Send("3", h.MainMarkup)
	})
}
