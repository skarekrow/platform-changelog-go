package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/redhatinsights/platform-changelog-go/internal/db"
)

func GetCommitsAll(w http.ResponseWriter, r *http.Request) {
	incRequests(r.URL.Path, r.Method, r.UserAgent())
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db.GetCommitsAll(db.DB))
}