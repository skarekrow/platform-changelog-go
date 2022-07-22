package db

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
	"gorm.io/gorm"
)

func GetDeploysAll(db *gorm.DB) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeploysAll)
	defer callDurationTimer.ObserveDuration()
	var deploys []structs.TimelinesData
	result := db.Model(models.Timelines{}).Order("Timestamp desc").Where("timelines.type = ?", "deploy").Scan(&deploys)
	return result, deploys
}

func GetDeploysByService(db *gorm.DB, service models.Services) (*gorm.DB, []structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeploysByService)
	defer callDurationTimer.ObserveDuration()
	var deploys []structs.TimelinesData
	result := db.Model(models.Timelines{}).Order("Timestamp desc").Where("timelines.service_id = ?", service.ID).Where("timelines.type = ?", "deploy").Scan(&deploys)
	return result, deploys
}

func GetDeployByRef(db *gorm.DB, ref string) (*gorm.DB, structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetDeployByRef)
	defer callDurationTimer.ObserveDuration()
	var deploy structs.TimelinesData
	result := db.Model(models.Timelines{}).Where("timelines.ref = ?", ref).Where("timelines.type = ?", "deploy").Scan(&deploy)
	return result, deploy
}
