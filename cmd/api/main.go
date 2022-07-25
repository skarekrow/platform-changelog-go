package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/endpoints"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
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
	sub := chi.NewRouter().With(metrics.ResponseMetricsMiddleware)

	// Mount the root of the api router on /api/v1
	r.Mount("/api/v1", sub)
	r.Get("/", lubdub)

	mr.Get("/", lubdub)
	mr.Get("/healthz", lubdub)
	mr.Handle("/metrics", promhttp.Handler())

	sub.Post("/github-webhook", endpoints.GithubWebhook)
	sub.Post("/gitlab-webhook", endpoints.GitlabWebhook)

	sub.Get("/services", endpoints.GetServicesAll)
	sub.Get("/timelines", endpoints.GetTimelinesAll)
	sub.Get("/commits", endpoints.GetCommitsAll)
	sub.Get("/deploys", endpoints.GetDeploysAll)

	sub.Get("/services/{service}", endpoints.GetServiceByName)
	sub.Get("/services/{service}/timelines", endpoints.GetTimelinesByService)
	sub.Get("/services/{service}/commits", endpoints.GetCommitsByService)
	sub.Get("/services/{service}/deploys", endpoints.GetDeploysByService)

	sub.Get("/timelines/{ref}", endpoints.GetTimelineByRef)
	sub.Get("/commits/{ref}", endpoints.GetCommitByRef)
	sub.Get("/deploys/{ref}", endpoints.GetDeployByRef)

	srv := http.Server{
		Addr:    ":" + cfg.PublicPort,
		Handler: r,
	}

	msrv := http.Server{
		Addr:    ":" + cfg.MetricsPort,
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
