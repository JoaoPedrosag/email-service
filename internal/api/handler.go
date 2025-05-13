package api

import (
	"net/http"
	"strconv"

	"github.com/JoaoPedrosag/email-service/internal/db"
	"github.com/JoaoPedrosag/email-service/internal/model"

	"github.com/gin-gonic/gin"
)

func CreateAuthorizedIP(c *gin.Context) {
    var input struct {
        IP string `json:"ip" binding:"required,ip"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	var existing model.AuthorizedIP
    if err := db.DB.Where("ip = ?", input.IP).First(&existing).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": "IP address is already authorized",
        })
        return
    }

    ip := model.AuthorizedIP{
        IP:      input.IP,
        Disable: true,
    }

    if err := db.DB.Create(&ip).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create IP"})
        return
    }

    c.JSON(http.StatusCreated, ip)
}

func ListAuthorizedIPs(c *gin.Context) {
    var ips []model.AuthorizedIP
    if err := db.DB.Find(&ips).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list IPs"})
        return
    }

    c.JSON(http.StatusOK, ips)
}

func ToggleAuthorizedIP(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    var ip model.AuthorizedIP
    if err := db.DB.First(&ip, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "IP not found"})
        return
    }

    ip.Disable = !ip.Disable

    if err := db.DB.Save(&ip).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update IP"})
        return
    }

    c.JSON(http.StatusOK, ip)
}
