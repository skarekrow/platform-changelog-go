package db

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
	"gorm.io/gorm"
)

func CreateCommitEntry(db *gorm.DB, t []models.Timelines) *gorm.DB {
	callDurationTimer := prometheus.NewTimer(metrics.SqlCreateCommitEntry)
	defer callDurationTimer.ObserveDuration()

	for _, timeline := range t {
		db.Create(&timeline)
	}

	return db
}

func GetCommitsAll(db *gorm.DB) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsAll)
	defer callDurationTimer.ObserveDuration()
	var commits []structs.TimelinesData
	result := db.Model(models.Timelines{}).Order("Timestamp desc").Where("timelines.type = ?", "commit").Scan(&commits)
	return result, commits
}

func GetCommitByRef(db *gorm.DB, ref string) (*gorm.DB, structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitByRef)
	defer callDurationTimer.ObserveDuration()
	var commit structs.TimelinesData
	result := db.Model(models.Timelines{}).Where("timelines.ref = ?", ref).Where("timelines.type = ?", "commit").Scan(&commit)
	return result, commit
}

func GetCommitsByService(db *gorm.DB, service structs.ServicesData) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsByService)
	defer callDurationTimer.ObserveDuration()
	var commits []structs.TimelinesData
	result := db.Model(models.Timelines{}).Where("timelines.service_id = ?", service.ID).Where("timelines.type = ?", "commit").Scan(&commits)
	return result, commits
}
