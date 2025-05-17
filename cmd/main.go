package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/JoaoPedrosag/email-service/internal/api"
	"github.com/JoaoPedrosag/email-service/internal/db"
	"github.com/JoaoPedrosag/email-service/internal/mailer"
	"github.com/JoaoPedrosag/email-service/internal/rabbitmq"
)

func main() {
	// Init DB
	db.Init()
	if err := db.RunMigrations(db.DB, "./migrations"); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Init Mailer
	mail := mailer.New()

    // Init Producer
	producer, err := rabbitmq.NewProducer()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer producer.Close()

	// Start Consumer
	go rabbitmq.StartEmailConsumer(mail)

	// Pass producer to API
	api.Init(producer)

	// Setup HTTP server
	r := gin.Default()

	r.POST("/authorized-ips", api.CreateAuthorizedIP)
	r.GET("/authorized-ips", api.ListAuthorizedIPs)
	r.PATCH("/authorized-ips/:id/toggle", api.ToggleAuthorizedIP)
	r.POST("/emails", api.EnqueueEmail)

	r.Run(":8080")
}
