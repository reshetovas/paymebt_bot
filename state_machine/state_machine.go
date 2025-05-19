package state_machine

import (
	"fmt"
	"payment_bot/models"
	"payment_bot/services"
	"payment_bot/utils"
	"strconv"
	"strings"
	"time"
)

type StateMachine struct {
	paymentService *services.PaymentService
	reportService  *services.ReportService
	exportService  *services.ExportService
}

func NewStateMachine(paymentService *services.PaymentService, reportService *services.ReportService, exportService *services.ExportService) *StateMachine {
	return &StateMachine{
		paymentService: paymentService,
		reportService:  reportService,
		exportService:  exportService,
	}
}

func (s *StateMachine) StateMachine(userID int64, text string) string {

	state, err := s.paymentService.GetUserState(userID)
	if err != nil {
		return "Ошибка при обработке запроса"
	}

	switch state {
	case models.AwaitingCategory:
		if err := s.paymentService.ProcessCategoryInput(userID, text); err != nil {
			return "Ошибка при сохранении категории"
		}
		return "Введи сумму:"

	case models.AwaitingAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return "Неверная сумма."
		}
		if err := s.paymentService.ProcessAmountInput(userID, amount); err != nil {
			return "Ошибка при сохранении суммы"
		}
		return "Укажи дату (ГГГГ-ММ-ДД) или 'сегодня':"

	case models.AwaitingDate:
		var dt time.Time
		if strings.ToLower(text) == "сегодня" {
			dt = time.Now()
		} else {
			dt, err = time.Parse("2006-01-02", text)
			if err != nil {
				return "Неверная дата."
			}
		}

		payment, err := s.paymentService.ProcessDateInput(userID, dt)
		if err != nil {
			return "Не удалось записать дать"
		}
		return fmt.Sprintf("✅ Добавлен платёж: %s, %.2f ₽, %s",
			payment.Category, payment.Amount, payment.Date.Format("02 Jan 2006"))
	case models.AwaitingCustomReport:
		from, to, err := utils.ParseCustomReportDates(text)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		report, err := s.reportService.GenerateCategoryReport(userID, from, to)
		if err != nil {
			return fmt.Sprintf("Ошибка при генерации отчёта")
		}
		return fmt.Sprintf(report)
	case models.AwaitingExport:
		from, to, err := utils.ParseCustomReportDates(text)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		report, err := s.exportService.ExportPayments(userID, from, to)
		if err != nil {
			return fmt.Sprintf("Ошибка при выгрузке платежей")
		}
		return fmt.Sprintf(report)
	default:
		// Обработка команд /report и /export...
		return ""
	}
}
