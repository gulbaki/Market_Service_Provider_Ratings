package api

import (
	"encoding/json"
	"log"
	"net/http"
	"notification_service/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	notiService service.NotificationService
}

func NewNotificationHandler(notiService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notiService: notiService}
}

// GetNotifications godoc
// @Summary      Get notifications by provider
// @Description  Retrieves new notifications for the specified provider and clears them afterwards
// @Tags         notifications
// @Produce      json
// @Param        providerId   path      int   true  "Provider ID"
// @Success      200  {array}  domain.Notification
// @Failure      400  {string}  string  "Invalid providerId"
// @Router       /notifications/{providerId} [get]
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	providerIDStr := vars["providerId"]
	providerID, err := strconv.Atoi(providerIDStr)
	if err != nil {
		http.Error(w, "Invalid providerId", http.StatusBadRequest)
		return
	}

	notifs := h.notiService.GetNotificationsByProviderID(providerID)
	if len(notifs) == 0 {
		// If there are no notifications, return an empty array
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
	log.Print(notifs)

	_ = h.notiService.ClearNotificationsByProviderID(providerID)
}
