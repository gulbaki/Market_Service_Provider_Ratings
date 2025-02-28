package service

import (
	"errors"
	"notification_service/internal/domain"
	"sync"
)

type NotificationService interface {
	CreateNotification(notification domain.Notification) error
	GetNotificationsByProviderID(providerID int) []domain.Notification
	ClearNotificationsByProviderID(providerID int) error
}

// InMemoryNotificationService
type InMemoryNotificationService struct {
	mu            sync.Mutex
	notifications map[int][]domain.Notification // providerId -> slice of notifications
}

func NewInMemoryNotificationService() *InMemoryNotificationService {
	return &InMemoryNotificationService{
		notifications: make(map[int][]domain.Notification),
	}
}

func (s *InMemoryNotificationService) CreateNotification(notification domain.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	providerID := notification.ProviderID
	s.notifications[providerID] = append(s.notifications[providerID], notification)

	return nil
}

func (s *InMemoryNotificationService) GetNotificationsByProviderID(providerID int) []domain.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	notifs, ok := s.notifications[providerID]
	if !ok {
		return []domain.Notification{} // boş
	}
	return notifs
}

func (s *InMemoryNotificationService) ClearNotificationsByProviderID(providerID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.notifications[providerID]
	if !ok {
		return errors.New("no notifications found for this provider")
	}
	delete(s.notifications, providerID)
	return nil
}
