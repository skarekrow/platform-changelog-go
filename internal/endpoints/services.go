package endpoints

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
)

func GetServicesAll(w http.ResponseWriter, r *http.Request) {
	incRequests(r.URL.Path, r.Method, r.UserAgent())
	start := time.Now()
	result, services := db.GetServicesAll(db.DB)
	elapsed := time.Since(start)
	observeDBTime("GetServicesAll", elapsed)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

func GetAllByServiceName(w http.ResponseWriter, r *http.Request) {
	incRequests(r.URL.Path, r.Method, r.UserAgent())
	serviceName := chi.URLParam(r, "service")
	start := time.Now()
	result, service := db.GetAllByServiceName(db.DB, serviceName)
	elapsed := time.Since(start)
	observeDBTime("GetAllByServiceName", elapsed)
	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Service not found"))
		return
	}
	l.Log.Debugf("URL Param: %s", serviceName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}