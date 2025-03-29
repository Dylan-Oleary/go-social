package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Dylan-Oleary/go-social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtxKey userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := getUserFromCtx(r)

	// TODO: Revert when using Auth
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), userToFollow.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	type status struct {
		OK bool `json:"ok"`
	}
	if err := app.jsonResponse(w, http.StatusOK, &status{OK: true}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userToUnfollow := getUserFromCtx(r)

	// TODO: Revert when using Auth
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, err)
		return
	}

	if err := app.store.Followers.Unfollow(r.Context(), userToUnfollow.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	type status struct {
		OK bool `json:"ok"`
	}

	if err := app.jsonResponse(w, http.StatusOK, &status{OK: true}); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userId)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundError(w, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	return r.Context().Value(userCtxKey).(*store.User)
}
