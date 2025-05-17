package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/JoaoPedrosag/email-service/internal/db"
	"github.com/JoaoPedrosag/email-service/internal/model"
)

func CreateAuthorizedIP(c *gin.Context) {
    var input struct {
        IP string `json:"ip" binding:"required,ip"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or missing IP address"})
        return
    }

    var exists int
    if err := db.DB.Get(&exists, "SELECT 1 FROM authorized_ips WHERE ip=$1", input.IP); err == nil {
        c.JSON(http.StatusConflict, gin.H{"message": "IP address is already authorized"})
        return
    }

    var lastID int64
    err := db.DB.QueryRow(
        `INSERT INTO authorized_ips (ip, disabled) VALUES ($1, $2) RETURNING id`,
        input.IP, true,
    ).Scan(&lastID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create authorized IP"})
        return
    }

    var ip model.AuthorizedIP
    if err := db.DB.Get(&ip,
        "SELECT id, ip, disabled, created_at, updated_at FROM authorized_ips WHERE id=$1",
        lastID,
    ); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve created IP record"})
        return
    }

    c.JSON(http.StatusCreated, ip)
}



func ListAuthorizedIPs(c *gin.Context) {
    ips := []model.AuthorizedIP{}
    if err := db.DB.Select(&ips, "SELECT * FROM authorized_ips"); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list authorized IPs"})
        return
    }
    c.JSON(http.StatusOK, ips)
}


func ToggleAuthorizedIP(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
        return
    }

    var ip model.AuthorizedIP
    if err := db.DB.Get(&ip, "SELECT * FROM authorized_ips WHERE id=$1", id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "IP not found"})
        return
    }

    newStatus := !ip.Disabled
    now := time.Now()
    if _, err := db.DB.Exec(
        "UPDATE authorized_ips SET disabled=$1, updated_at=$2 WHERE id=$3",
        newStatus, now, id,
    ); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update IP status"})
        return
    }

    ip.Disabled = newStatus
    ip.UpdatedAt = now
    c.JSON(http.StatusOK, ip)
}
