package endpoints

import (
	"net/http"
	"strconv"

	"github.com/redhatinsights/platform-changelog-go/internal/structs"
)

func writeResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func initQuery(r *http.Request) (structs.Query, error) {
	q := structs.Query{
		Page:  0,
		Limit: 10,
	}

	var err error

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page != "" {
		q.Page, err = strconv.Atoi(page)
	}

	if limit != "" {
		q.Limit, err = strconv.Atoi(limit)
	}

	return q, err
}
