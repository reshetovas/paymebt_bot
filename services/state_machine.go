package services

import (
	"fmt"
	"payment_bot/models"
	"strconv"
	"strings"
	"time"
)

func (s *PaymentService) StateMachine(userID int64, text string) string {

	state, err := s.GetUserState(userID)
	if err != nil {
		return "Ошибка при обработке запроса"
	}

	switch state {
	case models.AwaitingCategory:
		if err := s.ProcessCategoryInput(userID, text); err != nil {
			return "Ошибка при сохранении категории"
		}
		return "Введи сумму:"

	case models.AwaitingAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return "Неверная сумма."
		}
		if err := s.ProcessAmountInput(userID, amount); err != nil {
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

		payment, err := s.ProcessDateInput(userID, dt)
		if err != nil {
			return "Не удалось записать дать"
		}
		return fmt.Sprintf("✅ Добавлен платёж: %s, %.2f ₽, %s",
			payment.Category, payment.Amount, payment.Date.Format("02 Jan 2006"))
	case models.AwaitingCustomReport:
		from, to, err := s.ParseCustomReportDates(text)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		report, err := s.GenerateCategoryReport(userID, from, to)
		if err != nil {
			return fmt.Sprintf("Ошибка при генерации отчёта")
		}
		return fmt.Sprintf(report)
	case models.AwaitingExport:
		from, to, err := s.ParseCustomReportDates(text)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		report, err := s.ExportPayments(userID, from, to)
		if err != nil {
			return fmt.Sprintf("Ошибка при выгрузке платежей")
		}
		return fmt.Sprintf(report)
	default:
		// Обработка команд /report и /export...
		return ""
	}
}
