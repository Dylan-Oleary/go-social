package main

import (
	"log"
	"time"

	"github.com/Dylan-Oleary/go-social/internal/db"
	"github.com/Dylan-Oleary/go-social/internal/env"
	"github.com/Dylan-Oleary/go-social/internal/mailer"
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
		env:         env.GetString("ENV", "development"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
			mailTrap: mailTrapConfig{
				apiKey:    env.GetString("MAILTRAP_API_KEY", ""),
				fromEmail: env.GetString("MAILTRAP_FROM_EMAIL", ""),
			},
			sendGrid: sendGridConfig{
				apiKey:    env.GetString("SENDGRID_API_KEY", ""),
				fromEmail: env.GetString("SENDGRID_FROM_EMAIL", ""),
			},
		},
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

	// SendGrid
	// mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.sendGrid.fromEmail, logger)

	// MailTrap
	mailer, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.mailTrap.fromEmail)
	if err != nil {
		logger.Fatal("Failed to instantiate mailer")
	}

	app := &application{
		config: cfg,
		logger: logger,
		mailer: mailer,
		store:  store,
	}
	mux := app.mount()

	logger.Fatal(app.run(mux))
}
