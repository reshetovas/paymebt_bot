package handlers

import (
	"gopkg.in/telebot.v3"
	"payment_bot/services"
	"payment_bot/state_machine"
	"strings"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	reportService  *services.ReportService
	exportService  *services.ExportService
	stateMachine   *state_machine.StateMachine
}

func NewPaymentHandler(paymentService *services.PaymentService, reportService *services.ReportService, exportService *services.ExportService, stateMachine *state_machine.StateMachine) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		reportService:  reportService,
		exportService:  exportService,
		stateMachine:   stateMachine,
	}
}

func (h *PaymentHandler) HandleStart(c telebot.Context) error {
	return c.Send("Привет! Используй /add_payment, /report или /export.")
}

func (h *PaymentHandler) HandleAddPayment(c telebot.Context) error {
	userID := c.Sender().ID
	if err := h.paymentService.StartPaymentCreation(userID); err != nil {
		return c.Send("Ошибка при создании платежа")
	}
	markup := h.paymentService.GetMarkups()
	return c.Send("Выбери категорию:", markup)
}

func (h *PaymentHandler) HandleCategorySelection(c telebot.Context) error {
	userID := c.Sender().ID
	payload := strings.TrimPrefix(c.Callback().Data, "category:")

	err := h.paymentService.ProcessCategoryInput(userID, payload)
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
		report, err := h.reportService.GenerateCategoryReportToday(userID)
		if err != nil {
			return c.Send("Ошибка при генерации отчета.")
		}
		return c.Send(report)
	case "report_month":
		report, err := h.reportService.GenerateCategoryReportMonth(userID)
		if err != nil {
			return c.Send("Ошибка при генерации отчета.")
		}
		return c.Send(report)
	case "report_custom":
		userID := c.Sender().ID
		if err := h.reportService.StartCustomReportCreation(userID); err != nil {
			return c.Send("Ошибка при создании отчета")
		}
		return c.Send("Введи период в формате ГГГГ-ММ-ДД - ГГГГ-ММ-ДД")
	default:
		return nil
	}
}

func (h *PaymentHandler) HandleReport(c telebot.Context) error {
	markup := h.reportService.GetMarkupsReport()
	return c.Send("Выберите период отчета:", markup)
}

func (h *PaymentHandler) HandleExport(c telebot.Context) error {
	userID := c.Sender().ID
	if err := h.exportService.StartExportPayments(userID); err != nil {
		return c.Send("%s", err)
	}
	return c.Send("Введи период в формате ГГГГ-ММ-ДД - ГГГГ-ММ-ДД")
}

func (h *PaymentHandler) HandleTextInput(c telebot.Context) error {
	userID := c.Sender().ID
	text := strings.TrimSpace(c.Text())

	result := h.stateMachine.StateMachine(userID, text)
	return c.Send(result)
}
