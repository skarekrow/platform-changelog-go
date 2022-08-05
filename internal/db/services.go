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

func CreateServiceTableEntry(db *gorm.DB, name string, s config.Service) (result *gorm.DB, service models.Services) {
	newService := models.Services{Name: name, DisplayName: s.DisplayName, GHRepo: s.GHRepo, GLRepo: s.GLRepo, Branch: s.Branch, Namespace: s.Namespace, DeployFile: s.DeployFile}
	results := db.Create(&newService)
	return results, newService
}

func GetServicesAll(db *gorm.DB, offset int, limit int) (*gorm.DB, []structs.ExpandedServicesData, int64) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServicesAll)
	defer callDurationTimer.ObserveDuration()

	var count int64
	var services []structs.ExpandedServicesData

	dbQuery := db.Model(models.Services{})
	dbQuery.Find(&services).Count(&count)

	result := dbQuery.Limit(limit).Offset(offset).Find(&services)

	var servicesWithTimelines []structs.ExpandedServicesData
	for i := 0; i < len(services); i++ {
		_, _, s := GetLatest(db, services[i])

		servicesWithTimelines = append(servicesWithTimelines, s)
	}

	return result, servicesWithTimelines, count
}

func GetLatest(db *gorm.DB, service structs.ExpandedServicesData) (*gorm.DB, *gorm.DB, structs.ExpandedServicesData) {
	l.Log.Debugf("Query name: %s", service.Name)

	// TODO: Make one query to get the latest commit and deploy for each service
	comResult := db.Model(models.Timelines{}).Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", service.Name).Where("timelines.type = ?", "commit").Order("Timestamp desc").Limit(1).Find(&service.Commit)

	depResult := db.Model(models.Timelines{}).Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", service.Name).Where("timelines.type = ?", "deploy").Order("Timestamp desc").Limit(1).Find(&service.Deploy)

	return comResult, depResult, service
}

func GetServiceByName(db *gorm.DB, name string) (*gorm.DB, structs.ServicesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServiceByName)
	defer callDurationTimer.ObserveDuration()
	var service structs.ServicesData
	result := db.Model(models.Services{}).Where("name = ?", name).First(&service)
	return result, service
}

func GetServiceByGHRepo(db *gorm.DB, service_url string) (*gorm.DB, structs.ServicesData) {
	var service structs.ServicesData
	result := db.Model(models.Services{}).Where("gh_repo = ?", service_url).First(&service)

	return result, service
}
