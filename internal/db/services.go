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

func GetServicesAll(db *gorm.DB) (*gorm.DB, []structs.ServicesData) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServicesAll)
	defer callDurationTimer.ObserveDuration()
	var services []structs.ServicesData
	result := db.Model(models.Services{}).Find(&services)

	var servicesWithTimelines []structs.ServicesData
	for i := 0; i < len(services); i++ {
		_, s := GetLatest(db, services[i])
		servicesWithTimelines = append(servicesWithTimelines, s)
	}

	return result, servicesWithTimelines
}

func GetLatest(db *gorm.DB, service structs.ServicesData) (*gorm.DB, structs.ServicesData) {
	l.Log.Debugf("Query name: %s", service.Name)
	result := db.Model(models.Timelines{}).Select("*").Joins("JOIN services ON timelines.service_id = services.id").Where("services.name = ?", service.Name).Where("timelines.type = ?", "commit").Order("Timestamp desc").Limit(1).Find(&service.Timelines)
	return result, service
}

func GetServiceByName(db *gorm.DB, name string) (*gorm.DB, models.Services) {
	callDurationTimer := prometheus.NewTimer(metrics.SqlGetServiceByName)
	defer callDurationTimer.ObserveDuration()
	var service models.Services
	result := db.Where("name = ?", name).First(&service)
	return result, service
}

func GetServiceByGHRepo(db *gorm.DB, service_url string) (*gorm.DB, models.Services) {
	var service models.Services
	result := db.Where("gh_repo = ?", service_url).First(&service)

	return result, service
}
