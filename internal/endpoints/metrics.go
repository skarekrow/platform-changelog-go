package endpoints

import (
	"net/http"
	"strconv"
	"time"

	p "github.com/prometheus/client_golang/prometheus"
	pa "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requests = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_requests_total",
		Help: "Total number of requests",
	}, []string{})

	ghWebhooks = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_github_webhooks_total",
		Help: "Total number of github webhooks",
	}, []string{})

	glWebhooks = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_gitlab_webhooks_total",
		Help: "Total number of gitlab webhooks",
	}, []string{})

	ghWebhooksErrors = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_github_webhooks_errors_total",
		Help: "Total number of github webhooks errors",
	}, []string{})

	glWebhooksErrors = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_gitlab_webhooks_errors_total",
		Help: "Total number of gitlab webhooks errors",
	}, []string{})

	responseCodes = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_response_codes_total",
		Help: "Total number of response codes",
	}, []string{"code"})

	dbElapsed = pa.NewHistogramVec(p.HistogramOpts{
		Name:    "platform_changelog_db_elapsed_seconds",
		Help:    "Elapsed time for database operations",
	}, []string{"operation"})
)

type metricsTrackingResponseWriter struct {
	Wrapped http.ResponseWriter
	UserAgent string
}

func incRequests() {
	requests.With(p.Labels{}).Inc()
}

func observeDBTime(elapsed time.Duration) {
	dbElapsed.With(p.Labels{}).Observe(elapsed.Seconds())
}

func (m *metricsTrackingResponseWriter) Header() http.Header {
	return m.Wrapped.Header()
}

func (m *metricsTrackingResponseWriter) WriteHeader(statusCode int) {
	responseCodes.With(p.Labels{"code": strconv.Itoa(statusCode)}).Inc()
	m.Wrapped.WriteHeader(statusCode)
}

func (m *metricsTrackingResponseWriter) Write(b []byte) (int, error) {
	return m.Wrapped.Write(b)
}

func ResponseMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := &metricsTrackingResponseWriter{
			Wrapped: w,
			UserAgent: r.Header.Get("User-Agent"),
		}
		next.ServeHTTP(mw, r)
	})
}