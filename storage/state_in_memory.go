package storage

import (
	"fmt"
	"payment_bot/models"
	"sync"

	"github.com/rs/zerolog/log"
)

type InMemeoryStateStorage struct {
	UserStates  map[int64]models.State
	UserContext map[int64]*models.Payment
	mutex       sync.RWMutex
}

func NewStateStorage() *InMemeoryStateStorage {
	return &InMemeoryStateStorage{
		UserStates:  make(map[int64]models.State),
		UserContext: make(map[int64]*models.Payment),
	}
}

type StateStorage interface {
	GetUserState(userId int64) (models.State, error)
	UploadUserState(userId int64, state models.State) error
	DeleteUserState(userId int64) error
	GetUserContext(userId int64) (*models.Payment, error)
	UploadUserContext(userId int64, updateFunc func(*models.Payment)) error
	DeleteUserContext(userId int64) error
}

func (s *InMemeoryStateStorage) GetUserState(userId int64) (models.State, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	state, ok := s.UserStates[userId]
	if !ok {
		log.Error().Msgf("user %d state not found", userId)
		return "", fmt.Errorf("user %d state not found", userId)
	}
	return state, nil
}

func (s *InMemeoryStateStorage) UploadUserState(userId int64, state models.State) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.UserStates[userId] = state
	if s.UserStates[userId] == "" {
		log.Error().Msgf("user %d, state %v not set", userId, state)
		return fmt.Errorf("user %d, state %v not set", userId, state)
	}
	return nil
}

func (s *InMemeoryStateStorage) DeleteUserState(userId int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.UserStates, userId)
	if s.UserStates[userId] != "" {
		log.Error().Msgf("user %d faild to delete entry", userId)
		return fmt.Errorf("user %d faild to delete entry", userId)
	}

	return nil
}

func (s *InMemeoryStateStorage) GetUserContext(userId int64) (*models.Payment, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	payment, ok := s.UserContext[userId]
	if !ok {
		log.Error().Msgf("user %d state not found", userId)
		return nil, fmt.Errorf("user %d state not found", userId)
	}

	paymentCopy := *payment
	return &paymentCopy, nil
}

func (s *InMemeoryStateStorage) UploadUserContext(userId int64, updateFunc func(*models.Payment)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	payment, ok := s.UserContext[userId]
	if !ok {
		log.Error().Msgf("payment %v not found", payment)
		return fmt.Errorf("payment %v not found", payment)
	}

	updateFunc(payment)
	return nil
}

func (s *InMemeoryStateStorage) DeleteUserContext(userId int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.UserContext, userId)
	if s.UserStates[userId] != "" {
		log.Error().Msgf("faild to delete entry userID: %d", userId)
		return fmt.Errorf("faild to delete entry userID: %d", userId)
	}

	return nil
}
