package main

import (
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: Auth User Id From Token
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(43))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, 200, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
