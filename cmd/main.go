package main

import (
	"github.com/gin-gonic/gin"

	"github.com/JoaoPedrosag/email-service/internal/api"
	"github.com/JoaoPedrosag/email-service/internal/db"
	"github.com/JoaoPedrosag/email-service/internal/kafka"
)

func main() {
    db.Init()
    kafka.InitProducer()
    go kafka.StartEmailConsumer()

    r := gin.Default()

    r.POST("/authorized-ips", api.CreateAuthorizedIP)
    r.GET("/authorized-ips", api.ListAuthorizedIPs)
    r.PATCH("/authorized-ips/:id/toggle", api.ToggleAuthorizedIP)
    r.POST("/emails", api.EnqueueEmail)

    r.Run(":8080") 
}
