package endpoints

import (
	"net/http"

	"github.com/google/go-github/github"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
)

func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	
	payload, err := github.ValidatePayload(r, []byte(config.Get().GithubWebhookSecretKey))
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
	case *github.PullRequest:
		// not remotely complete. This just a placeholder.
		if e.GetMerged() {
			writeResponse(w, http.StatusOK, `{"msg": "merged"}`)
			return
		}
		writeResponse(w, http.StatusOK, `{"msg": "not merged yet"}`)
		return
	default:
		writeResponse(w, http.StatusOK, `{"msg": "Event from this repo is not a push event"}`)
		return	
	}
}