package db

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
	"gorm.io/gorm"
)

func GetServiceByName(db *gorm.DB, name string) (*gorm.DB, models.Services) {
	var service models.Services
	result := db.Where("name = ?", name).First(&service)
	return result, service
}

func CreateServiceTableEntry(db *gorm.DB, name string, s config.Service) (result *gorm.DB, service models.Services) {
	newService := models.Services{Name: name, DisplayName: s.DisplayName, GHRepo: s.GHRepo, GLRepo: s.GLRepo, Branch: s.Branch, Namespace: s.Namespace, DeployFile: s.DeployFile}
	results := db.Create(&newService)
	return results, newService
}

func GetServiceByGHRepo(db *gorm.DB, service_url string) (*gorm.DB, models.Services) {
	var service models.Services
	result := db.Where("gh_repo = ?", service_url).First(&service)

	return result, service
}

func CreateCommitEntry(db *gorm.DB, c []models.Commits) *gorm.DB {
	callDurationTimer := prometheus.NewTimer(metrics.SqlCreateCommitEntry)
	defer callDurationTimer.ObserveDuration()

	for _, commit := range c {
		db.Create(&commit)

		// Add a timeline entry for this commit to the db
		timeline := models.Timelines{CommitID: commit.ID, ServiceID: commit.ServiceID, Timestamp: commit.Timestamp}
		db.Omit("DeployID", "Deploy").Create(&timeline)
	}

	return db
}

// GetAllByServiceName returns all commits for a given service
// TODO: this should include deploys once we have support for that
func GetAllByServiceName(db *gorm.DB, name string) (*gorm.DB, models.Services) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetAllByServiceName)
	defer callDurationTimer.ObserveDuration()
	var service models.Services

	l.Log.Debugf("Query name: %s", name)

	// TODO: this should be one query that utilizes the structs defined in structs.go
	db.Table("services").Select("*").Where("name = ?", name).First(&service)
	result := db.Table("commits").Select("*").Joins("JOIN services ON commits.service_id = services.id").Where("services.name = ?", name).Order("Timestamp desc").Find(&service.Commits)
	service.Deploys = []models.Deploys{}
	return result, service
}

func GetLatest(db *gorm.DB, service models.Services) (*gorm.DB, models.Services) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetAllByServiceName)
	defer callDurationTimer.ObserveDuration()
	l.Log.Debugf("Query name: %s", service.Name)
	result := db.Table("commits").Select("*").Joins("JOIN services ON commits.service_id = services.id").Where("services.name = ?", service.Name).Order("Timestamp desc").Limit(1).Find(&service.Commits)
	return result, service
}

func GetServicesAll(db *gorm.DB) (*gorm.DB, []models.Services, []models.Services) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServicesAll)
	defer callDurationTimer.ObserveDuration()
	var services []models.Services
	result := db.Find(&services)

	var servicesWithCommits []models.Services
	for i := 0; i < len(services); i++ {
		_, s := GetLatest(db, services[i])
		servicesWithCommits = append(servicesWithCommits, s)
	}

	return result, services, servicesWithCommits
}

func GetCommitsAll(db *gorm.DB) (*gorm.DB, []models.Commits) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsAll)
	defer callDurationTimer.ObserveDuration()
	var commits []models.Commits
	result := db.Order("Timestamp desc").Find(&commits)
	return result, commits
}

func GetDeploysAll(db *gorm.DB) (*gorm.DB, []models.Deploys) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeploysAll)
	defer callDurationTimer.ObserveDuration()
	var deploys []models.Deploys
	result := db.Order("Timestamp desc").Find(&deploys)
	return result, deploys
}

/**
 * GetTimeline returns a timeline of commits and deploys for a service
 */
func GetTimeline(db *gorm.DB, service models.Services) (*gorm.DB, []structs.TimelineData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimeline)
	defer callDurationTimer.ObserveDuration()

	var timeline []structs.TimelineData

	// Joins the timeline table to the commits and deploys tables and into the TimelineData struct
	result := db.Debug().Model(models.Timelines{}).Select("*").Joins("JOIN commits ON timelines.commit_id = commits.id").Where("timelines.service_id = ?", service.ID).Scan(&timeline)

	return result, timeline
}
