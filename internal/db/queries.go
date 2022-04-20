package db

import (
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	l "github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"gorm.io/gorm"
)

func GetServiceByName(db *gorm.DB, name string) (*gorm.DB, models.Services) {
	var service models.Services
	result := db.Where("name = ?", name).First(&service)
	return result, service
}

func CreateServiceTableEntry(db *gorm.DB, name string, s config.Service) (result *gorm.DB, service models.Services) {
	newService := models.Services{Name: name, DisplayName: s.DisplayName, GHRepo: s.GHRepo, Branch: s.Branch, Namespace: s.Namespace, DeployFile: s.DeployFile}
	results := db.Create(&newService)
	return results, newService
}

func GetServiceByGHRepo(db *gorm.DB, service_url string) (*gorm.DB, models.Services) {
	var service models.Services
	result := db.Where("gh_repo = ?", service_url).First(&service)
	return result, service
}

func CreateCommitEntry(db *gorm.DB, c []models.Commits) *gorm.DB {
	for _, commit := range c {
		db.Create(&commit)
	}
	return db
}

func GetServicesAll(db *gorm.DB) ([]models.Services) {
	var services []models.Services
	db.Find(&services)
	return services
}

func GetCommitsAll(db *gorm.DB) ([]models.Commits) {
	var commits []models.Commits
	db.Find(&commits)
	return commits
}

func GetDeploysAll(db *gorm.DB) ([]models.Deploys) {
	var deploys []models.Deploys
	db.Find(&deploys)
	return deploys
}

// GetAllByServiceName returns all commits for a given service
// TODO: this should include deploys once we have support for that
func GetAllByServiceName(db *gorm.DB, name string) (*gorm.DB, models.Services) {
	var services models.Services
	l.Log.Debugf("Query name: %s", name)
	db.Table("services").Select("*").Where("name = ?", name).First(&services)
	result := db.Table("commits").Select("*").Joins("JOIN services ON commits.service_id = services.id").Where("services.name = ?", name).Find(&services.Commits)
	return result, services
}