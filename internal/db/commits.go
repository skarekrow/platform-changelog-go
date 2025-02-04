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

func GetCommitsAll(db *gorm.DB, offset int, limit int) (*gorm.DB, []models.Timelines, int64) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsAll)
	defer callDurationTimer.ObserveDuration()

	var count int64
	var commits []models.Timelines

	db = db.Model(models.Timelines{}).Where("timelines.type = ?", "commit")

	db.Find(&commits).Count(&count)
	result := db.Order("Timestamp desc").Limit(limit).Offset(offset).Find(&commits)

	return result, commits, count
}

func GetCommitsByService(db *gorm.DB, service structs.ServicesData, offset int, limit int) (*gorm.DB, []models.Timelines, int64) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitsByService)
	defer callDurationTimer.ObserveDuration()

	var count int64
	var commits []models.Timelines

	db = db.Model(models.Timelines{}).Where("timelines.service_id = ?", service.ID).Where("timelines.type = ?", "commit")

	db.Find(&commits).Count(&count)
	result := db.Order("Timestamp desc").Limit(limit).Offset(offset).Find(&commits)

	return result, commits, count
}

func GetCommitByRef(db *gorm.DB, ref string) (*gorm.DB, models.Timelines) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetCommitByRef)
	defer callDurationTimer.ObserveDuration()
	var commit models.Timelines
	result := db.Model(models.Timelines{}).Where("timelines.ref = ?", ref).Where("timelines.type = ?", "commit").Scan(&commit)
	return result, commit
}
