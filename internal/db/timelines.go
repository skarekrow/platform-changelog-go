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
