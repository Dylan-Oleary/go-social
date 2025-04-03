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

// ActivateUser godoc
//
//	@Summary		Activates a user
//	@Description	Activates a user
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation Token"
//	@Success		204		{string}	string "User activated"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	ctx := r.Context()
	err := app.store.Users.Activate(ctx, token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
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

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
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
