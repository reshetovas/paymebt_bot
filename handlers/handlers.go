package handlers

import (
	"gopkg.in/telebot.v3"
	"payment_bot/services"
	"strings"
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
	markup := h.service.GetMarkups()
	return c.Send("Выбери категорию:", markup)
}

func (h *PaymentHandler) HandleCategorySelection(c telebot.Context) error {
	userID := c.Sender().ID
	payload := strings.TrimPrefix(c.Callback().Data, "category:")

	err := h.service.ProcessCategoryInput(userID, payload)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка записи категории"})
	}

	err = c.Respond() // Удалим крутилку
	if err != nil {
		return err
	}

	return c.Send("Введи сумму:")
}

func (h *PaymentHandler) HandleReportSelection(c telebot.Context) error {
	userID := c.Sender().ID
	payload := strings.TrimPrefix(c.Callback().Data, "category:")

	switch payload {
	case "report_today":
		report, err := h.service.GenerateCategoryReportToday(userID)
		if err != nil {
			return c.Send("Ошибка при генерации отчета.")
		}
		return c.Send(report)
	case "report_month":
		report, err := h.service.GenerateCategoryReportMonth(userID)
		if err != nil {
			return c.Send("Ошибка при генерации отчета.")
		}
		return c.Send(report)
	case "report_custom":
		userID := c.Sender().ID
		if err := h.service.StartCustomReportCreation(userID); err != nil {
			return c.Send("Ошибка при создании отчета")
		}
		return c.Send("Введи период в формате ГГГГ-ММ-ДД - ГГГГ-ММ-ДД")
	default:
		return nil
	}
}

func (h *PaymentHandler) HandleReport(c telebot.Context) error {
	markup := h.service.GetMarkupsReport()
	return c.Send("Выберите период отчета:", markup)
}

func (h *PaymentHandler) HandleExport(c telebot.Context) error {
	userID := c.Sender().ID
	if err := h.service.StartExportPayments(userID); err != nil {
		return c.Send("%s", err)
	}
	return c.Send("Введи период в формате ГГГГ-ММ-ДД - ГГГГ-ММ-ДД")
}

func (h *PaymentHandler) HandleTextInput(c telebot.Context) error {
	userID := c.Sender().ID
	text := strings.TrimSpace(c.Text())

	result := h.service.StateMachine(userID, text)
	return c.Send(result)
}
