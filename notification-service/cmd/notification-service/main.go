package main

import (
	"log"
	"net/http"
	_ "notification_service/cmd/notification-service/docs"
	"notification_service/internal/api"
	"notification_service/internal/consumer"
	"notification_service/internal/service"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found or can't load it. Proceeding with system env variables.")
	}
}

// @title          Notification Service API
// @version        1.0
// @description    This service consumes rating-created events from Kafka and provides notifications for providers.
//
// @contact.name   Your Name
// @contact.email  your-email@example.com
// @BasePath       /
func main() {
	// 1. Config
	kafkaBrokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "rating-created"
	}
	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "notification-group"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9191"
	}

	// 2. Notification Service (in-memory)
	notiService := service.NewInMemoryNotificationService()

	// 3. Kafka Consumer
	brokersList := []string{kafkaBrokers}
	ratingConsumer := consumer.NewRatingConsumer(brokersList, kafkaTopic, kafkaGroupID, notiService)
	ratingConsumer.Start()
	defer ratingConsumer.Stop()

	// 4. Router (Gorilla Mux)
	r := mux.NewRouter()

	notiHandler := api.NewNotificationHandler(notiService)
	r.HandleFunc("/notifications/{providerId}", notiHandler.GetNotifications).Methods("GET")

	// Swagger endpoint
	//  => /swagger/index.html
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Notification Service running on %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
