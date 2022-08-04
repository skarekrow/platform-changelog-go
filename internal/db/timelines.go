package db

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redhatinsights/platform-changelog-go/internal/metrics"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"github.com/redhatinsights/platform-changelog-go/internal/structs"
	"gorm.io/gorm"
)

/**
 * GetTimeline returns a timeline of commits and deploys for a service
 */
func GetTimelinesAll(db *gorm.DB, page int, limit int) (*gorm.DB, []structs.TimelinesData, int64) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimelinesAll)
	defer callDurationTimer.ObserveDuration()

	var count int64
	var timelines []structs.TimelinesData

	// Concatanate the timeline fields
	fields := fmt.Sprintf("%s,%s,%s", strings.Join(timelinesFields, ","), strings.Join(commitsFields, ","), strings.Join(deploysFields, ","))

	db = db.Model(models.Timelines{}).Select(fields)

	db.Find(&timelines).Count(&count)
	result := db.Order("Timestamp desc").Limit(limit).Offset(page * limit).Scan(&timelines)

	return result, timelines, count
}

func GetTimelinesByService(db *gorm.DB, service structs.ServicesData, page int, limit int) (*gorm.DB, []structs.TimelinesData, int64) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimelinesByService)
	defer callDurationTimer.ObserveDuration()

	var count int64
	var timelines []structs.TimelinesData

	// Concatanate the timeline fields
	fields := fmt.Sprintf("%s,%s,%s", strings.Join(timelinesFields, ","), strings.Join(commitsFields, ","), strings.Join(deploysFields, ","))

	db = db.Model(models.Timelines{}).Select(fields).Where("service_id = ?", service.ID)

	db.Find(&timelines).Count(&count)
	result := db.Order("Timestamp desc").Limit(limit).Offset(page * limit).Scan(&timelines)

	return result, timelines, count
}

func GetTimelineByRef(db *gorm.DB, ref string) (*gorm.DB, structs.TimelinesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetTimelineByRef)
	defer callDurationTimer.ObserveDuration()

	var timeline structs.TimelinesData

	result := db.Model(models.Timelines{}).Select("*").Where("timelines.ref = ?", ref).Scan(&timeline)

	return result, timeline
}
