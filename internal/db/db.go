package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Falha ao carregar .env: %v", err)
    }

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    log.Printf("DSN: host=%s user=%s password=%s dbname=%s port=%s\n",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PORT"),
)


    dbConn, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatalf("Falha ao conectar no PostgreSQL: %v", err)
    }

    DB = dbConn
    log.Println("Conectado ao PostgreSQL via sqlx.")
}
