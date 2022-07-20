package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"github.com/google/uuid"
	"strings"

	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
)

type CommitsJSON struct {
	Name string `json:"name"`
	ServiceID int `json:"service_id"`
	Repo	  string `json:"repo"`
	Ref	  string `json:"ref"`
	Timestamp time.Time `json:"timestamp"`
	Author	  string `json:"author"`
	MergedBy  string `json:"merged_by"`
	Message  string `json:"message"`
}

type DeploysJSON struct {
	Name string `json:"name"`
	ServiceID int `json:"service_id"`
	Ref string `json:"ref"`
	Namespace string `json:"namespace"`
	Cluster string `json:"cluster"`
	Image string `json:"image"`
	Timestamp time.Time `json:"timestamp"`
}

type Fields struct {
	Commits  []CommitsJSON `json:"commits"`
	Deploys  []DeploysJSON `json:"deploys"`
}

// generatea  UUID with no hyphens
func generateUUID() string {
	uuidWithHyphens := uuid.New()
	uuid := strings.Replace(uuidWithHyphens.String(), "-", "", -1)
	return uuid
}

func main() {
	logging.InitLogger()

	cfg := config.Get()

	db.DbConnect(cfg)

	jsonFile, err := os.Open("development/db-seeder/db_seed.json")
	if err != nil {
		logging.Log.Fatal(err)
	}
	logging.Log.Info("seed.json opened")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var fields Fields

	err = json.Unmarshal(byteValue, &fields)
	if err != nil {
		logging.Log.Error("Cannot unmarshal JSON for seeder: %s", err)
	}

	fmt.Println("Seeding Commits Table")
	for i := 0; i < len(fields.Commits); i++ {
		// Query the DB for the service ID based on the repo
		var service models.Services
		db.DB.Table("services").Where("name = ?", fields.Commits[i].Name).First(&service)

		commits := models.Commits{
			ServiceID: service.ID,
			Repo: fields.Commits[i].Repo,
			Timestamp: fields.Commits[i].Timestamp,
			Author: fields.Commits[i].Author,
			MergedBy: fields.Commits[i].MergedBy,
			Message: fields.Commits[i].Message,
		}
		commits.Ref = generateUUID()
		logging.Log.Info(commits)
		db.DB.Create(&commits)
	}

	fmt.Println("Seeding Deploys Table")
	for i := 0; i < len(fields.Deploys); i++ {
		// Query the DB for the service ID based on the repo
		var service models.Services
		db.DB.Table("services").Where("name = ?", fields.Deploys[i].Name).First(&service)

		deploys := models.Deploys{
			ServiceID: service.ID,
			Namespace: fields.Deploys[i].Namespace,
			Cluster: fields.Deploys[i].Cluster,
			Image: fields.Deploys[i].Image,
			Timestamp: fields.Deploys[i].Timestamp,
		}
		deploys.Ref = generateUUID()
		logging.Log.Info(deploys)
		db.DB.Create(&deploys)
	}
	
	fmt.Println("DB seeding complete")
}
