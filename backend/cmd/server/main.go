package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"calculator-backend/internal/db"
	"calculator-backend/internal/history"
	"calculator-backend/internal/server"
)

func main() {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://calculator:calculator@localhost:5432/calculator?sslmode=disable"
	}

	pool, err := db.NewPool(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := history.NewRepository(pool)
	router := server.NewRouter(repo)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
