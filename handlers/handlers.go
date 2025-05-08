package handlers

import (
	"payment_bot/models"
	"payment_bot/services"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(s *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

func (h *PaymentHandler) HandleStart(c telebot.Context) error {
	return c.Send("Привет! Используй /add_payment, /report или /export.")
}

func (h *PaymentHandler) HandleAddPayment(c telebot.Context) error {
	userID := c.Sender().ID
	if err := h.service.StartPaymentCreation(userID); err != nil {
		return c.Send("Ошибка при создании платежа")
	}
	return c.Send("Укажи категорию:")
}

func (h *PaymentHandler) HandleReport(c telebot.Context) error {
	return c.Send("Формат: /report месяц|год\nПример: /report январь")
}

func (h *PaymentHandler) HandleExport(c telebot.Context) error {
	return c.Send("Формат: /export 2024-01-01 2024-03-31 [категория]")
}

func (h *PaymentHandler) HandleTextInput(c telebot.Context) error {
	userID := c.Sender().ID
	text := strings.TrimSpace(c.Text())

	state, err := h.service.GetUserState(userID)
	if err != nil {
		return c.Send("Ошибка при обработке запроса")
	}

	switch state {
	case models.AwaitingCategory:
		if err := h.service.ProcessCategoryInput(userID, text); err != nil {
			return c.Send("Ошибка при сохранении категории")
		}
		return c.Send("Введи сумму:")

	case models.AwaitingAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return c.Send("Неверная сумма.")
		}
		if err := h.service.ProcessAmountInput(userID, amount); err != nil {
			return c.Send("Ошибка при сохранении суммы")
		}
		return c.Send("Укажи дату (ГГГГ-ММ-ДД) или 'сегодня':")

	case models.AwaitingDate:
		var dt time.Time
		if text == "сегодня" {
			dt = time.Now()
		} else {
			dt, err = time.Parse("2006-01-02", text)
			if err != nil {
				return c.Send("Неверная дата.")
			}
		}

		if err := h.service.ProcessDateInput(userID, dt); err != nil {
			return c.Send("Не удалось записать дать")
		}
		return c.Send("Укажи дату (ГГГГ-ММ-ДД) или 'сегодня':")
	default:
		// Обработка команд /report и /export...
		return nil
	}
}
