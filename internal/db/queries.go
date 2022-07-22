package db

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
	"gorm.io/gorm"
)

var (
	timelinesFields = []string{"timelines.id", "timelines.timestamp", "timelines.service_id", "timelines.ref", "timelines.repo", "timelines.type"}
	commitsFields   = []string{"timelines.author", "timelines.message", "timelines.merged_by"}
	deploysFields   = []string{"timelines.namespace", "timelines.cluster", "timelines.image"}
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

func CreateCommitEntry(db *gorm.DB, t []models.Timelines) *gorm.DB {
	callDurationTimer := prometheus.NewTimer(metrics.SqlCreateCommitEntry)
	defer callDurationTimer.ObserveDuration()

	for _, timeline := range t {
		db.Create(&timeline)
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

	db.Table("services").Select("*").Where("name = ?", name).First(&service)
	result := db.Table("timelines").Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", name).Order("Timestamp desc").Find(&service.Timelines)

	return result, service
}

func GetLatest(db *gorm.DB, service models.Services) (*gorm.DB, models.Services) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetAllByServiceName)
	defer callDurationTimer.ObserveDuration()
	l.Log.Debugf("Query name: %s", service.Name)
	result := db.Table("timelines").Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", service.Name).Where("timelines.type = ?", "commit").Order("Timestamp desc").Limit(1).Find(&service.Timelines)
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

func GetCommitsAll(db *gorm.DB) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsAll)
	defer callDurationTimer.ObserveDuration()
	var commits []structs.TimelinesData
	result := db.Model(models.Timelines{}).Order("Timestamp desc").Where("timelines.type = ?", "commit").Scan(&commits)
	return result, commits
}

func GetDeploysAll(db *gorm.DB) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeploysAll)
	defer callDurationTimer.ObserveDuration()
	var deploys []structs.TimelinesData
	result := db.Model(models.Timelines{}).Order("Timestamp desc").Where("timelines.type = ?", "deploy").Scan(&deploys)
	return result, deploys
}

func GetDeployByRef(db *gorm.DB, ref string) (*gorm.DB, structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeployByRef)
	defer callDurationTimer.ObserveDuration()
	var deploy structs.TimelinesData
	result := db.Model(models.Timelines{}).Where("timelines.ref = ?", ref).Where("timelines.type = ?", "deploy").Scan(&deploy)
	return result, deploy
}

func GetCommitByRef(db *gorm.DB, ref string) (*gorm.DB, structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitByRef)
	defer callDurationTimer.ObserveDuration()
	var commit structs.TimelinesData
	result := db.Model(models.Timelines{}).Where("timelines.ref = ?", ref).Where("timelines.type = ?", "commit").Scan(&commit)
	return result, commit
}

func GetCommitsByService(db *gorm.DB, service models.Services) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsByService)
	defer callDurationTimer.ObserveDuration()
	var commits []structs.TimelinesData
	result := db.Model(models.Timelines{}).Where("timelines.service_id = ?", service.ID).Where("timelines.type = ?", "commit").Scan(&commits)
	return result, commits
}

func GetDeploysByService(db *gorm.DB, service models.Services) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeploysByService)
	defer callDurationTimer.ObserveDuration()
	var deploys []structs.TimelinesData
	result := db.Model(models.Timelines{}).Order("Timestamp desc").Where("timelines.service_id = ?", service.ID).Where("timelines.type = ?", "deploy").Scan(&deploys)
	return result, deploys
}

/**
 * GetTimeline returns a timeline of commits and deploys for a service
 */
func GetTimelinesAll(db *gorm.DB) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimelinesAll)
	defer callDurationTimer.ObserveDuration()

	var timelines []structs.TimelinesData

	// Concatanate the timeline fields
	fields := fmt.Sprintf("%s,%s,%s", strings.Join(timelinesFields, ","), strings.Join(commitsFields, ","), strings.Join(deploysFields, ","))

	// Joins the timeline table to the commits and deploys tables and into the TimelineData struct
	result := db.Model(models.Timelines{}).Select(fields).Scan(&timelines)

	return result, timelines
}

func GetTimelinesByService(db *gorm.DB, service models.Services) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimelinesByService)
	defer callDurationTimer.ObserveDuration()

	var timelines []structs.TimelinesData

	// Concatanate the timeline fields
	fields := fmt.Sprintf("%s,%s,%s", strings.Join(timelinesFields, ","), strings.Join(commitsFields, ","), strings.Join(deploysFields, ","))

	result := db.Model(models.Timelines{}).Select(fields).Where("service_id = ?", service.ID).Scan(&timelines)

	return result, timelines
}

func GetTimelineByRef(db *gorm.DB, ref string) (*gorm.DB, structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimelineByRef)
	defer callDurationTimer.ObserveDuration()

	var timeline structs.TimelinesData

	result := db.Model(models.Timelines{}).Select("*").Where("timelines.ref = ?", ref).Scan(&timeline)

	return result, timeline
}
