package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
)

func lubdub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("lubdub"))
}


func main() {

	logging.InitLogger()

	cfg := config.Get()

	db.DbConnect(cfg)

	r := chi.NewRouter()
	mr := chi.NewRouter()
	sub := chi.NewRouter()

	// Mount the root of the api router on /api/v1
	r.Mount("/api/v1", sub)
	r.Get("/", lubdub)

	mr.Get("/", lubdub)
	mr.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr: ":" + cfg.PublicPort,
		Handler: r,
	}

	msrv := http.Server{
		Addr: ":" + cfg.MetricsPort,
		Handler: mr,
	}

	go func() {
		if err := msrv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}