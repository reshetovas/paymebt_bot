package main

import (
	"fmt"
	"payment_bot/models"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

type State string

const (
	Idle             State = "idle"
	AwaitingCategory       = "awaiting_category"
	AwaitingAmount         = "awaiting_amount"
	AwaitingDate           = "awaiting_date"
)

var userStates = map[int64]State{}
var userContext = map[int64]*models.Payment{}

func RegisterHandlers(bot *telebot.Bot) {
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Привет! Используй /add_payment, /report или /export.")
	})

	bot.Handle("/add_payment", func(c telebot.Context) error {
		userID := c.Sender().ID
		userStates[userID] = AwaitingCategory
		userContext[userID] = &Payment{UserID: userID}
		return c.Send("Укажи категорию:")
	})

	bot.Handle("/report", func(c telebot.Context) error {
		return c.Send("Формат: /report месяц|год\nПример: /report январь")
	})

	bot.Handle("/export", func(c telebot.Context) error {
		return c.Send("Формат: /export 2024-01-01 2024-03-31 [категория]")
	})

	bot.Handle(telebot.OnText, handleFSMInput)
}

func handleFSMInput(c telebot.Context) error {
	userID := c.Sender().ID
	state := userStates[userID]
	text := strings.TrimSpace(c.Text())

	switch state {
	case AwaitingCategory:
		userContext[userID].Category = text
		userStates[userID] = AwaitingAmount
		return c.Send("Введи сумму:")

	case AwaitingAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return c.Send("Неверная сумма.")
		}
		userContext[userID].Amount = amount
		userStates[userID] = AwaitingDate
		return c.Send("Укажи дату (ГГГГ-ММ-ДД) или 'сегодня':")

	case AwaitingDate:
		var dt time.Time
		var err error
		if text == "сегодня" {
			dt = time.Now()
		} else {
			dt, err = time.Parse("2006-01-02", text)
			if err != nil {
				return c.Send("Неверная дата.")
			}
		}
		payment := userContext[userID]
		payment.Date = dt
		DB.Create(payment)

		// Очистка состояния
		delete(userStates, userID)
		delete(userContext, userID)

		return c.Send(fmt.Sprintf("✅ Добавлен платёж: %s, %.2f ₽, %s",
			payment.Category, payment.Amount, payment.Date.Format("02 Jan 2006")))

	default:
		text := c.Text()
		if strings.HasPrefix(text, "/report") {
			return generateReport(c)
		} else if strings.HasPrefix(text, "/export") {
			return exportPayments(c)
		}
		return c.Send("Напиши /add_payment для записи платежа.")
	}
}
