package kafka

import (
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var EmailWriter *kafka.Writer

func InitProducer() {
    broker := os.Getenv("KAFKA_BROKER") 
    topic  := os.Getenv("KAFKA_TOPIC")  

    EmailWriter = kafka.NewWriter(kafka.WriterConfig{
        Brokers:  []string{broker},
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
        Async:    false,
    })

    log.Printf("Kafka producer inicializado em %s (tópico %s)\n", broker, topic)
}

func SendEmailMessage(key string, value []byte) error {
    msg := kafka.Message{
        Key:   []byte(key),
        Value: value,
        Time:  time.Now(),
    }
    return EmailWriter.WriteMessages(nil, msg)
}