package handlers

import (
	"fmt"

	"github.com/cepmap/party-finder-bot/internal/config"
	"github.com/cepmap/party-finder-bot/internal/dbconnector"
	tele "gopkg.in/telebot.v4"
)

// Константы для текста кнопок
const (
	BtnSignups     = "Мои записи на события"
	BtnGames       = "Мои игры"
	BtnEvents      = "Мои события"
	BtnAdmin       = "Администрирование"
	BtnCreateEvent = "Создать событие"
	BtnBack        = "Назад"
)

type Handlers struct {
	db  *dbconnector.DBConnector
	cfg *config.Config
}

// Структура для хранения обработчиков кнопок
type buttonHandler struct {
	message   string
	getMarkup func(*dbconnector.User) *tele.ReplyMarkup
}

func NewHandlers(db *dbconnector.DBConnector, cfg *config.Config) *Handlers {
	return &Handlers{
		db:  db,
		cfg: cfg,
	}
}

// createMainMarkup создает главное меню с правильной структурой кнопок
func (h *Handlers) createMainMarkup(additionalButtons ...string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	// Базовые кнопки
	signupsBtn := markup.Text(BtnSignups)
	gamesBtn := markup.Text(BtnGames)
	eventsBtn := markup.Text(BtnEvents)

	rows := []tele.Row{
		markup.Row(signupsBtn),
		markup.Row(gamesBtn, eventsBtn),
	}

	// Добавляем дополнительные кнопки
	for _, text := range additionalButtons {
		btn := markup.Text(text)
		rows = append(rows, markup.Row(btn))
	}

	markup.Reply(rows...)
	return markup
}

// getMainMarkup создает главное меню с учетом роли пользователя
func (h *Handlers) getMainMarkup(user *dbconnector.User) *tele.ReplyMarkup {
	if user != nil && user.IsAdmin {
		return h.createMainMarkup(BtnAdmin)
	}
	return h.createMainMarkup()
}

// getEventsMarkup создает меню событий
func (h *Handlers) getEventsMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}
	createEventBtn := markup.Text(BtnCreateEvent)
	backBtn := markup.Text(BtnBack)

	markup.Reply(
		markup.Row(createEventBtn),
		markup.Row(backBtn),
	)

	return markup
}

func (h *Handlers) Register(b *tele.Bot) {
	b.Handle("/start", h.handleStart)
	b.Handle(tele.OnText, h.handleText)
}

// handleStart обрабатывает команду /start
func (h *Handlers) handleStart(c tele.Context) error {
	user, err := h.db.GetOrCreateUser(c.Sender().ID, h.cfg.AdminTelegramIDs)
	if err != nil {
		return c.Send("Ошибка при работе с базой данных")
	}

	markup := h.getMainMarkup(user)
	resStr := "Ваши ближайшие игры:"
	if !user.IsAdmin {
		resStr = "Вы не подписались ни на одну из игр"
	}

	return c.Send(fmt.Sprintf("Привет, %s \n%s", c.Sender().Username, resStr), markup)
}

// handleText обрабатывает нажатия на кнопки
func (h *Handlers) handleText(c tele.Context) error {
	text := c.Text()

	// Получаем пользователя из БД
	user, err := h.db.GetOrCreateUser(c.Sender().ID, h.cfg.AdminTelegramIDs)
	if err != nil {
		return c.Send("Ошибка при работе с базой данных")
	}

	// Map обработчиков кнопок
	handlers := map[string]buttonHandler{
		BtnSignups: {
			message:   "Вы записались на следующие события:",
			getMarkup: func(u *dbconnector.User) *tele.ReplyMarkup { return h.getMainMarkup(u) },
		},
		BtnGames: {
			message:   "Ваши игры:",
			getMarkup: func(u *dbconnector.User) *tele.ReplyMarkup { return h.getMainMarkup(u) },
		},
		BtnEvents: {
			message:   "Ваши события:",
			getMarkup: func(u *dbconnector.User) *tele.ReplyMarkup { return h.getEventsMarkup() },
		},
		BtnCreateEvent: {
			message:   "Создание события...",
			getMarkup: func(u *dbconnector.User) *tele.ReplyMarkup { return h.getEventsMarkup() },
		},
		BtnBack: {
			message:   "Главное меню",
			getMarkup: func(u *dbconnector.User) *tele.ReplyMarkup { return h.getMainMarkup(u) },
		},
		BtnAdmin: {
			message:   "Панель администратора",
			getMarkup: func(u *dbconnector.User) *tele.ReplyMarkup { return h.getMainMarkup(u) },
		},
	}

	// Выполняем обработчик
	if handler, exists := handlers[text]; exists {
		return c.Send(handler.message, handler.getMarkup(user))
	}

	return nil
}
