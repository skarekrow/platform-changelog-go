package endpoints

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	m "github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
	"github.com/redhatinsights/platform-changelog-go/internal/utils"
)

// GithubWebhook gets data from the webhook and enters it into the DB
func GithubWebhook(w http.ResponseWriter, r *http.Request) {

	var err error
	var payload []byte

	metrics.IncWebhooks("github", r.Method, r.UserAgent(), false)

	services := config.Get().Services

	if config.Get().Debug {
		payload, err = ioutil.ReadAll(r.Body)
	} else {
		payload, err = github.ValidatePayload(r, []byte(config.Get().GithubWebhookSecretKey))
	}
	if err != nil {
		l.Log.Error(err)
		metrics.IncWebhooks("github", r.Method, r.UserAgent(), true)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		l.Log.Error("could not parse webhook: err=%s\n", err)
		metrics.IncWebhooks("github", r.Method, r.UserAgent(), true)
		return
	}

	switch e := event.(type) {
	case *github.PingEvent:
		writeResponse(w, http.StatusOK, `{"msg": "ok"}`)
		return
	case *github.PushEvent:
		for key, service := range services {
			if service.GHRepo == e.Repo.GetURL() {
				_, s := db.GetServiceByName(db.DB, key)
				if s.Branch != strings.Split(utils.DerefString(e.Ref), "/")[2] {
					l.Log.Info("Branch mismatch: ", s.Branch, " != ", strings.Split(utils.DerefString(e.Ref), "/")[2])
					writeResponse(w, http.StatusOK, `{"msg": "Not a monitored branch"}`)
					return
				}
				commitData := getCommitData(e, s)
				result := db.CreateCommitEntry(db.DB, commitData)
				if result.Error != nil {
					l.Log.Errorf("Failed to insert webhook data: %v", result.Error)
					metrics.IncWebhooks("github", r.Method, r.UserAgent(), true)
					writeResponse(w, http.StatusInternalServerError, `{"msg": "Failed to insert webhook data"}`)
					return
				}
				l.Log.Infof("Created %d commit entries for %s", len(commitData), key)
				writeResponse(w, http.StatusOK, `{"msg": "ok"}`)
				return
			}
		}
		// catch for if the service is not registered
		l.Log.Infof("Service not found for %s", e.Repo.GetURL())
		writeResponse(w, http.StatusOK, `{"msg": "The service is not registered"}`)
		return
	default:
		l.Log.Errorf("Event type %T not supported", e)
		writeResponse(w, http.StatusOK, `{"msg": "Event from this repo is not a push event"}`)
		return
	}
}

func getCommitData(g *github.PushEvent, s structs.ServicesData) []m.Timelines {
	var commits []m.Timelines
	for _, commit := range g.Commits {
		record := m.Timelines{
			ServiceID: s.ID,
			Repo:      utils.DerefString(g.GetRepo().Name),
			Ref:       commit.GetID(),
			Type:      "commit",
			Timestamp: commit.Timestamp.Time,
			Author:    utils.DerefString(commit.GetAuthor().Login),
			MergedBy:  g.Pusher.GetName(),
			Message:   commit.GetMessage(),
		}
		commits = append(commits, record)
	}

	return commits
}
