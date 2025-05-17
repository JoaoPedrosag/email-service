package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/JoaoPedrosag/email-service/internal/db"
)


func Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Name	 string `json:"name" binding:"required"`

	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	var exists bool
	_ = db.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", body.Email)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3)",
		body.Email, body.Name, string(hash),
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var hash string
	err := db.DB.Get(&hash, "SELECT password_hash FROM users WHERE active = true AND email=$1", body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials or user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": body.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"email": body.Email,
	})
}
