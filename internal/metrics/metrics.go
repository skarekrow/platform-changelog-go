package metrics

import (
	"net/http"
	"strconv"

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

	SqlGetCommitsAll = pa.NewHistogram(p.HistogramOpts{
		Name:    "platform_changelog_sql_get_commits_all_seconds",
		Help:    "Elapsed time for sql lookup of all commits",
	})

	SqlGetServicesAll = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_services_all_seconds",
		Help: "Elapsed time for sql lookup of all services",
	})

	SqlGetAllByServiceName = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_all_by_service_name_seconds",
		Help: "Elapsed time for sql lookup of services by name",
	})

	SqlGetDeploysAll = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_deploys_all_seconds",
		Help: "Elapsed time for sql lookup of all deploys",
	})

	SqlCreateCommitEntry = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_create_commit_entry_seconds",
		Help: "Elapsed time for sql creation of commit entry",
	})
)

type MetricsTrackingResponseWriter struct {
	Wrapped http.ResponseWriter
	UserAgent string
}

func IncRequests(path string, method string, userAgent string) {
	requests.With(p.Labels{"path": path, "method": method, "user_agent": userAgent}).Inc()
}

func IncWebhooks(source string, method string, userAgent string, err bool) {
	if !err {
		webhooks.With(p.Labels{"source": source, "method": method, "user_agent": userAgent}).Inc()
	} else {
		webhookErrors.With(p.Labels{"source": source, "method": method, "user_agent": userAgent}).Inc()
	}
}

func (m *MetricsTrackingResponseWriter) Header() http.Header {
	return m.Wrapped.Header()
}

func (m *MetricsTrackingResponseWriter) WriteHeader(statusCode int) {
	responseCodes.With(p.Labels{"code": strconv.Itoa(statusCode)}).Inc()
	m.Wrapped.WriteHeader(statusCode)
}

func (m *MetricsTrackingResponseWriter) Write(b []byte) (int, error) {
	return m.Wrapped.Write(b)
}

func ResponseMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := &MetricsTrackingResponseWriter{
			Wrapped: w,
			UserAgent: r.Header.Get("User-Agent"),
		}
		next.ServeHTTP(mw, r)
	})
}