package api

import (
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
