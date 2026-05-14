package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

// 1. Using a Pointer Receiver (*PaginatedFeedQuery)
// 2. Returning a single error
func (fq *PaginatedFeedQuery) Parse(r *http.Request) error {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fmt.Errorf("invalid limit parameter: %w", err)
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fmt.Errorf("invalid offset parameter: %w", err)
		}
		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		parsedTime, err := parseTime(since)
		if err != nil {
			// 3. Stop immediately and return the error!
			return fmt.Errorf("invalid 'since' date format: %w", err)
		}
		fq.Since = parsedTime
	}

	until := qs.Get("until")
	if until != "" {
		parsedTime, err := parseTime(until)
		if err != nil {
			// 3. Stop immediately and return the error!
			return fmt.Errorf("invalid 'until' date format: %w", err)
		}
		fq.Until = parsedTime
	}

	return nil
}

func parseTime(s string) (string, error) {
	// 1. Try the full DateTime format (2026-05-01 15:30:00)
	t, err := time.Parse(time.DateTime, s)
	if err == nil {
		return t.Format(time.DateTime), nil
	}

	// 2. Fallback: Try the short Date format (2026-05-01)
	t, err = time.Parse(time.DateOnly, s)
	if err == nil {
		// We still format it as DateTime so PostgreSQL is happy!
		return t.Format(time.DateTime), nil
	}

	// 3. If it is garbage like "2027" or "apple", actually return the error
	return "", fmt.Errorf("invalid time format: %s", s)
}
