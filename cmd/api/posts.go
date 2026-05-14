package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type postKey string

const postctx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// @Summary		Create a post
// @Description	Creates a new post with a title, content, and optional tags.
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			payload	body		CreatePostPayload	true	"Post payload"
// @Success		201		{object}	store.Post
// @Failure		400		{object}	error	"Bad Request / Validation Error"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := getUserFromContext(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID, // TODO: Replace with authenticated user ID
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Get a post
// @Description	Fetches a post by its ID, including its associated comments.
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			postID	path		int	true	"Post ID"
// @Success		200		{object}	store.Post
// @Failure		404		{object}	error	"Post Not Found"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// @Summary		Delete a post
// @Description	Deletes a post by its ID.
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			postID	path	int	true	"Post ID"
// @Success		204		"No Content - Successfully deleted"
// @Failure		404		{object}	error	"Post Not Found"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")

	id, err := strconv.ParseInt(idParam, 10, 64)
	fmt.Println("id", id)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err = app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// @Summary		Update a post
// @Description	Partially updates a post's title or content by its ID.
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			postID	path		int					true	"Post ID"
// @Param			payload	body		UpdatePostPayload	true	"Post update payload"
// @Success		200		{object}	store.Post
// @Failure		400		{object}	error	"Bad Request / Validation Error"
// @Failure		404		{object}	error	"Post Not Found"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/posts/{postID} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// --- Middleware & Helpers (No Swagger comments needed here) ---

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")

		id, err := strconv.ParseInt(idParam, 10, 64)
		fmt.Println("id", id)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)

			}
			return
		}

		ctx = context.WithValue(ctx, postctx, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postctx).(*store.Post)

	return post
}
