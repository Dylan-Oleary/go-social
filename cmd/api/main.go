package main

import (
	"log"

	"github.com/Dylan-Oleary/go-social/internal/env"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{addr: env.GetString("ADDR", ":8080")}
	app := &application{config: cfg}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
