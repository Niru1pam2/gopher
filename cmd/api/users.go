package main

import (
	"context"
	"errors"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// @Summary		Get user profile
// @Description	Fetches a user's public profile by their ID.
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path		int	true	"User ID"
// @Success		200		{object}	store.User
// @Failure		400		{object}	error	"Bad Request / Invalid ID"
// @Failure		404		{object}	error	"User Not Found"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID < 1 {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.getUser(ctx, userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// @Summary		Follow a user
// @Description	Allows the authenticated user to follow another user.
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path	int			true	"ID of the user to follow"
// @Param			payload	body	FollowUser	true	"Follower User ID (Temporary Auth Bypass)"
// @Success		204		"No Content - Successfully followed"
// @Failure		400		{object}	error	"Bad Request / Malformed JSON"
// @Failure		404		{object}	error	"Target User Not Found"
// @Failure		409		{object}	error	"Conflict - Already following"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/v1/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Follower.Follow(ctx, followedUser.ID, followedID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Unfollow a user
// @Description	Allows the authenticated user to unfollow another user.
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path	int			true	"ID of the user to unfollow"
// @Param			payload	body	FollowUser	true	"Unfollower User ID (Temporary Auth Bypass)"
// @Success		204		"No Content - Successfully unfollowed"
// @Failure		400		{object}	error	"Bad Request / Malformed JSON"
// @Failure		404		{object}	error	"Target User Not Found"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/v1/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)
	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Follower.Unfollow(ctx, followerUser.ID, unfollowedID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

// --- Middleware & Helpers (No Swagger comments needed here) ---

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userId)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.badRequestError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
