package main

import (
	"net/http"
	"social/internal/store"
)

// @Summary		Get user feed
// @Description	Retrieves a paginated, searchable, and filterable feed of posts for the authenticated user.
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			limit	query		int		false	"Number of items to return (max 20)"	default(20)
// @Param			offset	query		int		false	"Number of items to skip"				default(0)
// @Param			sort	query		string	false	"Sorting order (asc or desc)"			default(desc)
// @Param			tags	query		string	false	"Comma-separated list of tags to filter by"
// @Param			search	query		string	false	"Text to search in title or content"
// @Param			since	query		string	false	"Start date for filtering (YYYY-MM-DD)"
// @Param			until	query		string	false	"End date for filtering (YYYY-MM-DD)"
// @Success		200		{array}		store.PostWithMetadata
// @Failure		400		{object}	error	"Bad Request / Invalid Parameters"
// @Failure		500		{object}	error	"Internal Server Error"
// @Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := &store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	// TODO: Replace '2' with the actual authenticated user's ID
	feed, err := app.store.Posts.GetUserFeed(ctx, 2, fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
