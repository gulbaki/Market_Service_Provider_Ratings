package consumer

import (
	"context"
	"encoding/json"
	"log"
	"notification_service/internal/domain"
	"notification_service/internal/service"
	"notification_service/pkg/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

type RatingCreatedEvent struct {
	ProviderId int       `json:"providerId"`
	Score      int       `json:"score"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
}

type RatingConsumer struct {
	reader      *kafka.Reader
	notiService service.NotificationService // Servis katmanına kaydedeceğiz
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewRatingConsumer constructor
func NewRatingConsumer(brokers []string, topic string, groupID string, notiService service.NotificationService) *RatingConsumer {
	ctx, cancel := context.WithCancel(context.Background())

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID, // Consumer Group
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &RatingConsumer{
		reader:      r,
		notiService: notiService,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (rc *RatingConsumer) Start() {
	go func() {
		defer rc.reader.Close()
		for {
			m, err := rc.reader.ReadMessage(rc.ctx)
			if err != nil {
				if err == context.Canceled {
					log.Println("Consumer context canceled. Stopping consumer.")
					return
				}
				log.Printf("Error reading message: %v\n", err)
				continue
			}

			var event RatingCreatedEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v\n", err)
				continue
			}

			n := domain.Notification{
				ProviderID: event.ProviderId,
				Score:      event.Score,
				Comment:    event.Comment,
				CreatedAt:  event.CreatedAt,
			}
			err = rc.notiService.CreateNotification(n)
			if err != nil {
				logger.Error("Failed to create notification:", err)
			} else {

				logger.WithFields(map[string]interface{}{
					"providerId": event.ProviderId,
					"score":      event.Score,
				}).Info("Simulating send of notification (received rating-created event).")
			}
		}
	}()
}

// Stop consumer (context cancel)
func (rc *RatingConsumer) Stop() {
	rc.cancel()
}
