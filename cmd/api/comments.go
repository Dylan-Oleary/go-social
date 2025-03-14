package main

import (
	"net/http"

	"github.com/Dylan-Oleary/go-social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) addCommentsToPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, err)
		return
	}

	post := getPostFromCtx(r)
	comment := store.Comment{
		Content: payload.Content,
		PostID:  post.ID,
		// TODO: Change after auth
		UserID: 1,
	}

	if err := app.store.Comments.Create(r.Context(), &comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, comment)
}
