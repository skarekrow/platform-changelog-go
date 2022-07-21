package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
)

/**
 * GetTimeline returns a timeline of commits for a given service
 */
func GetTimeline(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())

	serviceName := chi.URLParam(r, "service")

	result, service := db.GetServiceByName(db.DB, serviceName)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't find the service"))
		return
	}

	_, timeline := db.GetTimeline(db.DB, service)

	// if result.Error != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Error producing the timeline"))
	// 	w.Write([]byte(result.Error.Error()))
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timeline)
}
