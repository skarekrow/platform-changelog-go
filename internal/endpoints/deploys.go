package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
)

func GetDeploysAll(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())

	result, deploys := db.GetDeploysAll(db.DB)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deploys)
}

func GetDeploysByService(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())
	serviceName := chi.URLParam(r, "service")

	result, service := db.GetServiceByName(db.DB, serviceName)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Couldn't find the service"))
		return
	}

	result, deploys := db.GetDeploysByService(db.DB, service)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deploys)
}

func GetDeployByRef(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests(r.URL.Path, r.Method, r.UserAgent())
	ref := chi.URLParam(r, "ref")

	result, deploy := db.GetDeployByRef(db.DB, ref)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Deploy not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deploy)
}
