package main

import (
	"log"

	"github.com/Dylan-Oleary/go-social/internal/db"
	"github.com/Dylan-Oleary/go-social/internal/env"
	"github.com/Dylan-Oleary/go-social/internal/store"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5434/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store)
}
