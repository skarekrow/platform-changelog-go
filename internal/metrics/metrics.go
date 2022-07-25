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

	SqlCreateCommitEntry = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_create_commit_entry_seconds",
		Help: "Elapsed time for sql creation of commit entry",
	})

	SqlGetServicesAll = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_services_all_seconds",
		Help: "Elapsed time for sql lookup of all services",
	})

	SqlGetTimelinesAll = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_timelines_all_seconds",
		Help: "Elapsed time for sql lookup of timeline entries",
	})

	SqlGetCommitsAll = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_commits_all_seconds",
		Help: "Elapsed time for sql lookup of all commits",
	})

	SqlGetDeploysAll = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_deploys_all_seconds",
		Help: "Elapsed time for sql lookup of all deploys",
	})

	SqlGetServiceByName = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_service_by_name_seconds",
		Help: "Elapsed time for sql lookup of services by name",
	})

	SqlGetCommitsByService = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_commits_by_service_seconds",
		Help: "Elapsed time for sql lookup of commits by service",
	})

	SqlGetDeploysByService = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_deploys_by_service_seconds",
		Help: "Elapsed time for sql lookup of deploys by service",
	})

	SqlGetTimelinesByService = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_timelines_by_service_seconds",
		Help: "Elapsed time for sql lookup of a service's timeline entries",
	})

	SqlGetCommitByRef = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_commit_by_ref_seconds",
		Help: "Elapsed time for sql lookup of commit by ref",
	})

	SqlGetDeployByRef = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_deploy_by_ref_seconds",
		Help: "Elapsed time for sql lookup of deploy by ref",
	})

	SqlGetTimelineByRef = pa.NewHistogram(p.HistogramOpts{
		Name: "platform_changelog_sql_get_timeline_by_ref_seconds",
		Help: "Elapsed time for sql lookup of timeline by ref",
	})
)

type MetricsTrackingResponseWriter struct {
	Wrapped   http.ResponseWriter
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
			Wrapped:   w,
			UserAgent: r.Header.Get("User-Agent"),
		}
		next.ServeHTTP(mw, r)
	})
}
