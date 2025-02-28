package service_test

import (
	"testing"

	"notification_service/internal/domain"
	"notification_service/internal/service"
)

func TestInMemoryNotificationService_CreateAndGet(t *testing.T) {
	// Arrange
	notiService := service.NewInMemoryNotificationService()
	n := domain.Notification{
		ProviderID: 123,
		Score:      5,
		Comment:    "Great job",
	}

	// Act
	if err := notiService.CreateNotification(n); err != nil {
		t.Fatalf("failed to create notification: %v", err)
	}

	// Assert
	notifs := notiService.GetNotificationsByProviderID(123)
	if len(notifs) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifs))
	}

	got := notifs[0]
	if got.Score != 5 || got.Comment != "Great job" {
		t.Errorf("notification data mismatch: got %+v", got)
	}
}

func TestInMemoryNotificationService_ClearNotifications(t *testing.T) {
	// Arrange
	notiService := service.NewInMemoryNotificationService()
	n := domain.Notification{
		ProviderID: 101,
		Score:      4,
		Comment:    "Test comment",
	}
	_ = notiService.CreateNotification(n)

	// Act
	err := notiService.ClearNotificationsByProviderID(101)
	if err != nil {
		t.Fatalf("failed to clear notifications: %v", err)
	}

	// Assert
	notifs := notiService.GetNotificationsByProviderID(101)
	if len(notifs) != 0 {
		t.Errorf("expected 0 notifications after clear, got %d", len(notifs))
	}
}

func TestInMemoryNotificationService_ClearNotifications_NonExistingProvider(t *testing.T) {
	// Arrange
	notiService := service.NewInMemoryNotificationService()

	// Act
	err := notiService.ClearNotificationsByProviderID(999)

	// Assert
	if err == nil {
		t.Errorf("expected error when clearing notifications for non-existing provider")
	}
}
