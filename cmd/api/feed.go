package main

import (
	"net/http"

	"github.com/Dylan-Oleary/go-social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginationFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "asc",
		Tags:   []string{},
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, err)
		return
	}

	ctx := r.Context()

	// TODO: Auth User Id From Token
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(43), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, 200, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
