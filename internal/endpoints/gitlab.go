package endpoints

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"fmt"

	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	m "github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/xanzy/go-gitlab"
)

func GetURL(p *gitlab.PushEvent) string {
	if p == nil || p.Repository == nil {
		return ""
	}
	return p.Repository.URL
}

func GetRepo(p *gitlab.PushEvent) *gitlab.Repository {
	if p == nil || p.Repository == nil {
		return nil
	}
	return p.Repository
}

func GetID(p *struct {
	ID        string     "json:\"id\""
	Message   string     "json:\"message\""
	Title     string     "json:\"title\""
	Timestamp *time.Time "json:\"timestamp\""
	URL       string     "json:\"url\""
	Author    struct {
		Name  string "json:\"name\""
		Email string "json:\"email\""
	} "json:\"author\""
	Added    []string "json:\"added\""
	Modified []string "json:\"modified\""
	Removed  []string "json:\"removed\""
}) string {
	if p == nil {
		return ""
	}
	return p.ID
}

func GetTime(p *time.Time) time.Time {
	if p == nil {
		return time.Time{}
	}
	return *p
}

type PingEvent struct {
	// Random string of GitHub zen.
	Zen *string `json:"zen,omitempty"`
	// The ID of the webhook that triggered the ping.
	HookID *int64 `json:"hook_id,omitempty"`
	// The webhook configuration.
	Hook *gitlab.Hook `json:"hook,omitempty"`
	//Installation *Installation `json:"installation,omitempty"`
}

func GetAuthor(p *struct {
	ID        string     "json:\"id\""
	Message   string     "json:\"message\""
	Title     string     "json:\"title\""
	Timestamp *time.Time "json:\"timestamp\""
	URL       string     "json:\"url\""
	Author    struct {
		Name  string "json:\"name\""
		Email string "json:\"email\""
	} "json:\"author\""
	Added    []string "json:\"added\""
	Modified []string "json:\"modified\""
	Removed  []string "json:\"removed\""
}) string {
	if p == nil {
		return ""
	}
	return p.Author.Name
}

func GetName(p *gitlab.PushEvent) string {
	if p == nil {
		return ""
	}
	return p.UserName
}

func GetMessage(p *struct {
	ID        string     "json:\"id\""
	Message   string     "json:\"message\""
	Title     string     "json:\"title\""
	Timestamp *time.Time "json:\"timestamp\""
	URL       string     "json:\"url\""
	Author    struct {
		Name  string "json:\"name\""
		Email string "json:\"email\""
	} "json:\"author\""
	Added    []string "json:\"added\""
	Modified []string "json:\"modified\""
	Removed  []string "json:\"removed\""
}) string {
	if p == nil {
		return ""
	}
	return p.Message
}

// GitlabWebhook gets data from the webhook and enters it into the DB
func GitlabWebhook(w http.ResponseWriter, r *http.Request) {

	var err error
	var payload []byte

	metrics.IncWebhooks("gitlab", r.Method, r.UserAgent(), false)

	services := config.Get().Services

	if config.Get().Debug {
		payload, err = ioutil.ReadAll(r.Body)
	} else if (r.Header["X-gitlab-token"][0]) == config.Get().GitlabWebhookSecretKey {

		payload, err = ioutil.ReadAll(r.Body)

	}
	if err != nil {
		l.Log.Error(err)
		metrics.IncWebhooks("gitlab", r.Method, r.UserAgent(), true)
		return
	}
	defer r.Body.Close()

	event, err := gitlab.ParseWebhook(gitlab.WebhookEventType(r), payload)
	if err != nil {
		l.Log.Error("could not parse webhook: err=%s\n", err)
		metrics.IncWebhooks("gitlab", r.Method, r.UserAgent(), true)
		return
	}

	switch e := event.(type) {

	case PingEvent:
		writeResponse(w, http.StatusOK, `{"msg": "ok"}`)
		return

	case *gitlab.PushEvent:
		for key, service := range services {
			if service.GLRepo == GetURL(e) {
				_, s := db.GetServiceByName(db.DB, key)
				if s.Branch != strings.Split((e.Ref), "/")[2] {
					l.Log.Info("Branch mismatch: ", s.Branch, " != ", strings.Split((e.Ref), "/")[2])
					writeResponse(w, http.StatusOK, `{"msg": "Not a monitored branch"}`)
					return
				}
				commitData := getCommitData2(e, s)
				result := db.CreateCommitEntry(db.DB, commitData)
				if result.Error != nil {
					l.Log.Errorf("Failed to insert webhook data: %v", result.Error)
					metrics.IncWebhooks("gitlab", r.Method, r.UserAgent(), true)
					writeResponse(w, http.StatusInternalServerError, `{"msg": "Failed to insert webhook data"}`)
					return
				}
				db.DB.Commit()
				l.Log.Infof("Created %d commit entries for %s", len(commitData), key)
				writeResponse(w, http.StatusOK, `{"msg": "ok"}`)
				return
			}
		}
		// catch for if the service is not registered
		l.Log.Infof("Service not found for %s", GetURL(e))
		fmt.Println(GetURL(e))

		writeResponse(w, http.StatusOK, `{"msg": "The service is not registered"}`)
		return
	default:
		l.Log.Errorf("Event type %T not supported", e)
		writeResponse(w, http.StatusOK, `{"msg": "Event from this repo is not a push event"}`)
		return
	}
}

func getCommitData2(g *gitlab.PushEvent, s m.Services) []m.Commits {
	var commits []m.Commits
	for _, commit := range g.Commits {
		record := m.Commits{
			ServiceID: s.ID,
			Repo:      GetRepo(g).Name,
			Ref:       GetID(commit),
			Timestamp: GetTime(commit.Timestamp),
			Author:    GetAuthor(commit),
			MergedBy:  GetName(g),
			Message:   GetMessage(commit),
		}
		commits = append(commits, record)
	}

	return commits
}
