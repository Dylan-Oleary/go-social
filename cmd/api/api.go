package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dylan-Oleary/go-social/docs"
	"github.com/Dylan-Oleary/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	swagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	logger *zap.SugaredLogger
	store  store.Storage
}

type config struct {
	addr   string
	apiURL string
	db     dbConfig
	env    string
	mail   mailConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	exp time.Duration
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/swagger/*", swagger.Handler(swagger.URL(fmt.Sprintf("%s/swagger/doc.json", app.config.addr))))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.addCommentsToPostHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Route("/follow", func(r chi.Router) {
					r.Put("/", app.followUserHandler)
				})
				r.Route("/unfollow", func(r chi.Router) {
					r.Put("/", app.unfollowUserHandler)
				})
			})

			r.Group((func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			}))
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Post("/", app.registerUserHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
