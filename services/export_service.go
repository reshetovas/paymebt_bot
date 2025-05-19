package services

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"payment_bot/models"
	"payment_bot/storage"
	"strings"
	"time"
)

type ExportService struct {
	state   storage.StateStorage
	storage storage.PaymentStorage
}

func NewExportService(state storage.StateStorage, storage storage.PaymentStorage) *ExportService {
	return &ExportService{
		state:   state,
		storage: storage,
	}
}

func (s *ExportService) StartExportPayments(userID int64) error {
	if err := s.state.UploadUserState(userID, models.AwaitingExport); err != nil {
		log.Error().Err(err).Msg("Error uploading user state")
		return err
	}
	return nil
}

func (s *ExportService) ExportPayments(userID int64, from, to time.Time) (string, error) {
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, time.UTC)

	records, err := s.storage.GetPaymentsByPeriod(userID, from, to)
	if err != nil {
		log.Error().Err(err).Msg("Error getting payments by period")
		return "", err
	}

	fromDate := from.Format("2006-01-02")
	toDate := to.Format("2006-01-02")

	if len(records) == 0 {
		log.Error().Msgf("No payments found by period: %s - %s", fromDate, toDate)
		return "", fmt.Errorf("Нет платежей за период: %s - %s", fromDate, toDate)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Выгрузка платежей с %s по %s:\n\n", fromDate, toDate))

	for _, r := range records {
		DateDate := r.Date.Format("02 Jan 2006")
		sb.WriteString(fmt.Sprintf("%-15v %-15s %.2f ₽\n", DateDate, r.Category, r.Amount))
	}

	if err := s.state.DeleteUserState(userID); err != nil {
		log.Error().Err(err).Msg("Error deleting user state")
		return "", err
	}

	return sb.String(), nil
}
