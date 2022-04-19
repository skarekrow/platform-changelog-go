package endpoints

import (
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	m "github.com/redhatinsights/platform-changelog-go/internal/models"
)

// GithubWebhook gets data from the webhook and enters it into the DB
func GithubWebhook(w http.ResponseWriter, r *http.Request) {

	var err error
	var payload []byte 

	services := config.Get().Services

	if config.Get().Debug {
		payload, err = ioutil.ReadAll(r.Body)
	} else {
		payload, err = github.ValidatePayload(r, []byte(config.Get().GithubWebhookSecretKey))
	}
	if err != nil {
		l.Log.Error(err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		l.Log.Error("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PingEvent:
		writeResponse(w, http.StatusOK, `{"msg": "ok"}`)
		return
	case *github.PullRequestEvent:
		// not remotely complete. This just a placeholder.
		if *e.PullRequest.Merged {
			commit := getCommitData(e, services)
			result := db.CreateCommitEntry(db.DB, commit)
			db.DB.Commit()
			l.Log.Info("Created commit entry:", result.Statement)
			writeResponse(w, http.StatusOK, `{"msg": "merged"}`)
			return
		}
		writeResponse(w, http.StatusOK, `{"msg": "PR not merged yet"}`)
		return
	default:
		writeResponse(w, http.StatusOK, `{"msg": "Event from this repo is not a push event"}`)
		return
	}
}

func getCommitData(g *github.PullRequestEvent, s map[string]config.Service) m.Commits {
	commit := &m.Commits{
		Repo:      *g.Repo.Name,
		Ref:       *g.PullRequest.Head.Ref,
		Title:     *g.PullRequest.Title,
		Timestamp: *g.PullRequest.MergedAt,
		Author:    *g.PullRequest.GetUser().Login,
		MergedBy:  *g.PullRequest.MergedBy.Login,
		Message:   *g.PullRequest.Body,
	}

	return *commit
}
