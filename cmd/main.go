package main

import (
	"github.com/gin-gonic/gin"

	"github.com/JoaoPedrosag/email-service/internal/api"
	"github.com/JoaoPedrosag/email-service/internal/db"
	"github.com/JoaoPedrosag/email-service/internal/model"
)

func main() {
    db.Init()
    db.DB.AutoMigrate(&model.AuthorizedIP{})

    r := gin.Default()

    r.POST("/authorized-ips", api.CreateAuthorizedIP)
    r.GET("/authorized-ips", api.ListAuthorizedIPs)
    r.PATCH("/authorized-ips/:id/toggle", api.ToggleAuthorizedIP)

    r.Run(":8080") 
}
