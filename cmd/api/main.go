package main

import (
	"log"

	"github.com/Dylan-Oleary/go-social/internal/db"
	"github.com/Dylan-Oleary/go-social/internal/env"
	"github.com/Dylan-Oleary/go-social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Go Social
//	@description	social network written in Go
//	@termsOfService	http://swagger.io/terms/

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	err := env.LoadEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:adminpassword@localhost:5434/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewStorage(db)
	app := &application{config: cfg, logger: logger, store: store}
	mux := app.mount()

	logger.Fatal(app.run(mux))
}
