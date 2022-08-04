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
		Offset: 0,
		Limit:  10,
	}

	var err error

	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	if offset != "" {
		q.Offset, err = strconv.Atoi(offset)
	}

	if limit != "" {
		q.Limit, err = strconv.Atoi(limit)
	}

	return q, err
}
