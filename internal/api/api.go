package api

import (
	"github.com/JoaoPedrosag/email-service/internal/rabbitmq"
)

var producer *rabbitmq.Producer

func Init(p *rabbitmq.Producer) {
	producer = p
}
