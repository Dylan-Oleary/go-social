package main

import (
	"fmt"
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

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5434/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		fmt.Println(cfg.db)
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database connection pool established")

	store := store.NewStorage(db)
	app := &application{config: cfg, store: store}
	mux := app.mount()

	log.Fatal(app.run(mux))
}
