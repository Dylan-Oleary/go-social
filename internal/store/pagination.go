package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

const limitQsKey string = "limit"
const offsetQsKey string = "offset"
const sortQsKey string = "sort"
const tagsQsKey string = "tags"
const searchQsKey string = "search"
const sinceQsKey string = "since"
const untilQsKey string = "until"

func (fq PaginationFeedQuery) Parse(r *http.Request) (PaginationFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get(limitQsKey)
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}

		fq.Limit = l
	}

	offset := qs.Get(offsetQsKey)
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}

		fq.Offset = o
	}

	sort := qs.Get(sortQsKey)
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get(tagsQsKey)
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get(searchQsKey)
	if search != "" {
		fq.Search = search
	}

	since := qs.Get(sinceQsKey)
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := qs.Get(untilQsKey)
	if until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
