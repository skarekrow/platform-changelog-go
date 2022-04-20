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
	}, []string{"path", "method", "user_agent"})

	webhooks = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_webhooks_total",
		Help: "Total number of incoming webhooks",
	}, []string{"source", "method", "user_agent"})

	webhookErrors = pa.NewCounterVec(p.CounterOpts{
		Name: "platform_changelog_webhooks_errors_total",
		Help: "Total number of incoming webhook errors",
	}, []string{"source", "method", "user_agent"})

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

func incRequests(path string, method string, userAgent string) {
	requests.With(p.Labels{"path": path, "method": method, "user_agent": userAgent}).Inc()
}

func incWebhooks(source string, method string, userAgent string, err bool) {
	if !err {
		webhooks.With(p.Labels{"source": source, "method": method, "user_agent": userAgent}).Inc()
	} else {
		webhookErrors.With(p.Labels{"source": source, "method": method, "user_agent": userAgent}).Inc()
	}
}

func observeDBTime(operation string, elapsed time.Duration) {
	dbElapsed.With(p.Labels{"operation": operation}).Observe(elapsed.Seconds())
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