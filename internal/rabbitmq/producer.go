package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer wraps the RabbitMQ channel and queue info.
type Producer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewProducer() (*Producer, error) {
    amqpUser := os.Getenv("RABBITMQ_USER")
    amqpPass := os.Getenv("RABBITMQ_PASS")
    amqpHost := os.Getenv("RABBITMQ_HOST") // ex: "localhost:5672"
    queueName := os.Getenv("QUEUE_NAME")

    amqpURL := fmt.Sprintf("amqp://%s:%s@%s/", amqpUser, amqpPass, amqpHost)

    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to declare queue: %w", err)
    }

    log.Printf("RabbitMQ producer initialized on queue %s", queueName)

    return &Producer{
        conn:  conn,
        ch:    ch,
        queue: queue,
    }, nil
}

// Send serializes the given event and publishes it to the queue.
func (p *Producer) Send(evt any) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	err = p.ch.Publish(
		"",             // exchange
		p.queue.Name,   // routing key (queue name)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close gracefully closes the channel and connection.
func (p *Producer) Close() {
	if p.ch != nil {
		_ = p.ch.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
