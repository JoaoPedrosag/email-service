package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/JoaoPedrosag/email-service/internal/mailer"
	"github.com/JoaoPedrosag/email-service/internal/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartEmailConsumer(mail *mailer.Mailer) {
	rmqUser := os.Getenv("RABBITMQ_USER")
	rmqPass := os.Getenv("RABBITMQ_PASS")
	rmqHost := os.Getenv("RABBITMQ_HOST") // Ex: "localhost:5672"
	queueName := os.Getenv("QUEUE_NAME")

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/", rmqUser, rmqPass, rmqHost)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Printf("RabbitMQ consumer started on queue %s", queueName)

	for msg := range msgs {
		go processEmailMessage(mail, msg.Body)
	}
}

func processEmailMessage(mail *mailer.Mailer, body []byte) {
	var evt message.EmailEvent

	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("Invalid payload: %v\n", err)
		return
	}

	if err := mail.Send(evt.To, evt.Subject, evt.Body); err != nil {
		log.Printf("Failed to send email to %s: %v\n", evt.To, err)
	} else {
		log.Printf("Email successfully sent to %s\n", evt.To)
	}
}
