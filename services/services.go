// services/payment.go
package services

import (
	"payment_bot/models"
	"payment_bot/storage"
	"time"
)

type PaymentService struct {
	state   storage.StateStorage
	storage storage.PaymentStorage
}

func NewPaymentService(state storage.StateStorage, storage storage.PaymentStorage) *PaymentService {
	return &PaymentService{
		state:   state,
		storage: storage,
	}
}

func (s *PaymentService) StartPaymentCreation(userID int64) error {
	if err := s.state.UploadUserState(userID, models.AwaitingCategory); err != nil {
		return err
	}
	return s.state.UploadUserContext(userID, func(p *models.Payment) {
		p.UserID = userID
	})
}

func (s *PaymentService) GetUserState(userID int64) (models.State, error) {
	return s.state.GetUserState(userID)
}

func (s *PaymentService) GetPayment(userID int64) (*models.Payment, error) {
	return s.state.GetUserContext(userID)
}

func (s *PaymentService) ProcessCategoryInput(userID int64, category string) error {
	if err := s.state.UploadUserState(userID, models.AwaitingAmount); err != nil {
		return err
	}
	return s.state.UploadUserContext(userID, func(p *models.Payment) {
		p.Category = category
	})
}

func (s *PaymentService) ProcessAmountInput(userID int64, amount float64) error {
	if err := s.state.UploadUserState(userID, models.AwaitingDate); err != nil {
		return err
	}
	return s.state.UploadUserContext(userID, func(p *models.Payment) {
		p.Amount = amount
	})
}

func (s *PaymentService) ProcessDateInput(userID int64, date time.Time) error {
	//update payment info in memory
	err := s.state.UploadUserContext(userID, func(p *models.Payment) {
		p.Date = date
	})
	if err != nil {
		return err
	}

	// get payment info in memory
	payment, err := s.state.GetUserContext(userID)
	if err != nil {
		return err
	}

	// save payment info in storage
	if err := s.storage.SavePayment(payment); err != nil {
		return err
	}

	// clean memory
	if err := s.state.DeleteUserContext(userID); err != nil {
		return err
	}
	return s.state.DeleteUserContext(userID)
}

// func (s *PaymentService) GenerateReport(period string) (string, error) {
// 	// Логика формирования отчета
// }

// func (s *PaymentService) ExportPayments(from, to time.Time, category string) ([]*models.Payment, error) {
// 	// Логика экспорта
// }
