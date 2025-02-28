package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notification_service/internal/domain"
	"notification_service/internal/service"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestNotificationHandler_GetNotifications(t *testing.T) {

	notiService := service.NewInMemoryNotificationService()

	notiService.CreateNotification(domain.Notification{
		ProviderID: 101,
		Score:      5,
		Comment:    "Excellent service!",
	})

	handler := NewNotificationHandler(notiService)

	router := mux.NewRouter()
	router.HandleFunc("/notifications/{providerId}", handler.GetNotifications).Methods("GET")

	req, err := http.NewRequest(http.MethodGet, "/notifications/101", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}

	var notifs []domain.Notification
	if err := json.Unmarshal(rr.Body.Bytes(), &notifs); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(notifs) != 1 {
		t.Errorf("expected 1 notification, got %d", len(notifs))
	}

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected status 200 on second call, got %d", rr2.Code)
	}
	if !strings.Contains(rr2.Body.String(), "[]") {
		t.Errorf("expected empty array on second call, got %s", rr2.Body.String())
	}
}
