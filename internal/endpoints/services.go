package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
)

func GetServicesAll(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())
	//result, services := db.GetServicesAll(db.DB)

	result, _, latest_commits := db.GetServicesAll(db.DB)
	if result.Error != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(latest_commits)

}

func GetAllByServiceName(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())

	serviceName := chi.URLParam(r, "service")
	result, service := db.GetAllByServiceName(db.DB, serviceName)

	/**
	 * Couldn't get the commits table
	 */
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	/**
	 * the service doesn't exist
	 */
	if service.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Service not found"))
		return
	}
	l.Log.Debugf("URL Param: %s", serviceName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}
