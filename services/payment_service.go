// services/payment.go
package services

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"
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
		log.Error().Err(err).Msg("Error uploading user state")
		return err
	}
	return s.state.InitUserContext(userID, func(p *models.Payment) {
		p.UserID = userID
	})
}

func (s *PaymentService) GetMarkups() *telebot.ReplyMarkup {
	categories := []string{"Еда", "Транспорт", "Развлечения", "Прочее"}
	var buttons [][]telebot.InlineButton
	for _, cat := range categories {
		btn := telebot.InlineButton{
			Text: cat,
			Data: "category:" + cat,
		}
		buttons = append(buttons, []telebot.InlineButton{btn})
	}

	// Отправляем сообщение с кнопками
	markup := &telebot.ReplyMarkup{
		InlineKeyboard: buttons,
	}

	return markup
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

func (s *PaymentService) ProcessDateInput(userID int64, date time.Time) (*models.Payment, error) {
	//update payment info in memory
	err := s.state.UploadUserContext(userID, func(p *models.Payment) {
		p.Date = date
	})
	if err != nil {
		return nil, err
	}

	// get payment info in memory
	payment, err := s.state.GetUserContext(userID)
	if err != nil {
		return nil, err
	}

	// save payment info in storage
	if err := s.storage.SavePayment(payment); err != nil {
		return nil, err
	}

	// clean memory
	if err := s.state.DeleteUserContext(userID); err != nil {
		return nil, err
	}
	if err := s.state.DeleteUserState(userID); err != nil {
		return nil, err
	}
	return payment, nil
}
