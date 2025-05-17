package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/JoaoPedrosag/email-service/internal/message"
	"github.com/segmentio/kafka-go"
)

func StartEmailConsumer() {
    broker  := os.Getenv("KAFKA_BROKER")
    topic   := os.Getenv("KAFKA_TOPIC")
    groupID := os.Getenv("KAFKA_GROUP_ID")

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:        []string{broker},
        Topic:          topic,
        GroupID:        groupID,
        MinBytes:       10e3,
        MaxBytes:       10e6,
        CommitInterval: time.Second,
    })
    log.Printf("Kafka consumer iniciado no tópico %s (group %s)\n", topic, groupID)

    for {
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Erro lendo mensagem Kafka: %v\n", err)
            time.Sleep(time.Second * 5)
            continue
        }
        go processEmailMessage(m.Value)
    }
}

func processEmailMessage(value []byte) {
    var evt message.EmailEvent
    if err := json.Unmarshal(value, &evt); err != nil {
        log.Printf("Payload inválido: %v\n", err)
        return
    }

    // if err := mailer.Send(evt.To, evt.Subject, evt.Body); err != nil {
    //     log.Printf("Falha ao enviar e-mail para %s: %v\n", evt.To, err)
    // } else {
    //     log.Printf("E-mail enviado com sucesso para %s\n", evt.To)
    // }
}
