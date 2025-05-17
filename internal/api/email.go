package api

import (
	"log"
	"net/http"
	"time"

	"github.com/JoaoPedrosag/email-service/internal/message"
	"github.com/gin-gonic/gin"
)

func EnqueueEmail(c *gin.Context) {
	var evt message.EmailEvent
	if err := c.ShouldBindJSON(&evt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()
	if err := producer.Send(evt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to enqueue email"})
		return
	}
	log.Printf("Enqueue duration: %v", time.Since(start))

	c.JSON(http.StatusAccepted, gin.H{"status": "enqueued"})
}
