package api

import (
	"encoding/json"
	"net/http"

	"github.com/JoaoPedrosag/email-service/internal/kafka"
	"github.com/JoaoPedrosag/email-service/internal/message"
	"github.com/gin-gonic/gin"
)

func EnqueueEmail(c *gin.Context) {
    var evt message.EmailEvent
    if err := c.ShouldBindJSON(&evt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    data, err := json.Marshal(evt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao serializar mensagem"})
        return
    }

    if err := kafka.SendEmailMessage(evt.To, data); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao enfileirar e-mail"})
        return
    }

    c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})
}