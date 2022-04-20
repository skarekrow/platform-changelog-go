package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redhatinsights/platform-changelog-go/internal/db"
)

func GetDeploysAll(w http.ResponseWriter, r *http.Request) {
	incRequests(r.URL.Path, r.Method, r.UserAgent())
	start := time.Now()
	result, deploys := db.GetDeploysAll(db.DB)
	elapsed := time.Since(start)
	observeDBTime("GetDeploysAll", elapsed)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deploys)
}