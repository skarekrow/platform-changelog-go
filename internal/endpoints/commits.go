package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redhatinsights/platform-changelog-go/internal/db"
)

func GetCommitsAll(w http.ResponseWriter, r *http.Request) {
	incRequests(r.URL.Path, r.Method, r.UserAgent())
	start := time.Now()
	result, commits := db.GetCommitsAll(db.DB)
	elapsed := time.Since(start)
	observeDBTime("GetCommitsAll", elapsed)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commits)
}