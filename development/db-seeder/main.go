package main

import (
	"encoding/json"
	"time"
	"fmt"
	"os"
	"io/ioutil"

	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
)

type CommitsJSON struct {
	ServiceID int `json:"service_id"`
	Repo	  string `json:"repo"`
	Ref	  string `json:"ref"`
	Timestamp time.Time `json:"timestamp"`
	Author	  string `json:"author"`
	MergedBy  string `json:"merged_by"`
	Message  string `json:"message"`
}

type Fields struct {
	Services []models.Services `json:"services"`
	Commits  []CommitsJSON `json:"commits"`
}

func main() {
	cfg := config.Get()

	db.DbConnect(cfg)

	jsonFile, err := os.Open("db_seed.json")
	if err != nil {
		l.Log.Fatal(err)
	}
	l.Log.Info("seed.json opened")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var fields Fields

	json.Unmarshal(byteValue, &fields)

	fmt.Println("Seeding Services Table")
	for i := 0; i < len(fields.Services); i++ {
		l.Log.Info(fields.Services[i].Name)
		services := models.Services{
			Name: fields.Services[i].Name,
			DisplayName: fields.Services[i].DisplayName,
			GHRepo: fields.Services[i].GHRepo,
			GLRepo: fields.Services[i].GLRepo,
			Namespace: fields.Services[i].Namespace,
			Branch: fields.Services[i].Branch,
		}
		l.Log.Info(services)

		db.DB.Create(&services)
	}

	fmt.Println("Seeding Commits Table")
	for i := 0; i < len(fields.Commits); i++ {

		commits := models.Commits{
			ServiceID: fields.Commits[i].ServiceID,
			Repo: fields.Commits[i].Repo,
			Ref: fields.Commits[i].Ref,
			Timestamp: fields.Commits[i].Timestamp,
			Author: fields.Commits[i].Author,
			MergedBy: fields.Commits[i].MergedBy,
			Message: fields.Commits[i].Message,
		}
		l.Log.Info(commits)
		db.DB.Create(&commits)
	}
	
	fmt.Println("DB seeding complete")
}
