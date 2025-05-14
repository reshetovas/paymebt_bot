package services

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"
	"payment_bot/models"
	"strings"
	"time"
)

func (s *PaymentService) GetMarkupsReport() *telebot.ReplyMarkup {
	todayBtn := telebot.InlineButton{
		Data: "report_today",
		Text: "За сегодня",
	}
	monthBtn := telebot.InlineButton{
		Data: "report_month",
		Text: "За месяц",
	}
	customBtn := telebot.InlineButton{
		Data: "report_custom",
		Text: "Выбрать даты",
	}

	markup := &telebot.ReplyMarkup{}
	markup.InlineKeyboard = [][]telebot.InlineButton{
		{todayBtn},
		{monthBtn},
		{customBtn},
	}

	return markup
}

func (s *PaymentService) StartCustomReportCreation(userID int64) error {
	if err := s.state.UploadUserState(userID, models.AwaitingCustomReport); err != nil {
		log.Error().Err(err).Msg("Error uploading user state")
		return err
	}
	return nil
}

func (s *PaymentService) GenerateCategoryReportToday(userID int64) (string, error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, time.UTC)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	records, err := s.storage.GetCountByCategory(userID, from, to)
	if err != nil {
		return "", err
	}

	today := from.Format("2006-01-02")

	if len(records) == 0 {
		return "", fmt.Errorf("no payments found by period: %s - %s", today, today)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Отчёт по категориям с %s по %s:\n\n", today, today))

	for _, r := range records {
		sb.WriteString(fmt.Sprintf("%-15s %.2f ₽\n", r.Category, r.Amount))
	}

	return sb.String(), nil
}

func (s *PaymentService) GenerateCategoryReportMonth(userID int64) (string, error) {
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	log.Debug().Msgf("from: %v, to: %v", from, to)
	records, err := s.storage.GetCountByCategory(userID, from, to)
	if err != nil {
		return "", err
	}

	fromDate := from.Format("2006-01-02")
	toDate := to.Format("2006-01-02")

	if len(records) == 0 {
		return "", fmt.Errorf("no payments found by period: %s - %s", fromDate, toDate)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Отчёт по категориям с %s по %s:\n\n", fromDate, toDate))

	for _, r := range records {
		sb.WriteString(fmt.Sprintf("%-15s %.2f ₽\n", r.Category, r.Amount))
	}

	return sb.String(), nil
}

func (s *PaymentService) GenerateCategoryReport(userID int64, from, to time.Time) (string, error) {
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, time.UTC)
	records, err := s.storage.GetCountByCategory(userID, from, to)
	if err != nil {
		return "", err
	}

	fromDate := from.Format("2006-01-02")
	toDate := to.Format("2006-01-02")

	if len(records) == 0 {
		return "", fmt.Errorf("no payments found by period: %s - %s", fromDate, toDate)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Отчёт по категориям с %s по %s:\n\n", fromDate, toDate))

	for _, r := range records {
		sb.WriteString(fmt.Sprintf("%-15s %.2f ₽\n", r.Category, r.Amount))
	}

	if err := s.state.DeleteUserState(userID); err != nil {
		return "", err
	}

	return sb.String(), nil
}
