package db

import (
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
	"gorm.io/gorm"
)

func GetServiceByName(db *gorm.DB, service_id string) (*gorm.DB, models.Services) {
	var service models.Services
	result := db.Where("display_name = ?", service_id).First(&service)
	return result, service
}

func CreateServiceTableEntry(db *gorm.DB, s config.Service) (result *gorm.DB, service models.Services) {
	newService := models.Services{DisplayName: s.DisplayName, GHRepo: s.GHRepo, Branch: s.Branch, Namespace: s.Namespace, DeployFile: s.DeployFile}
	results := db.Create(&newService)
	return results, newService
}

func CreateCommitEntry(db *gorm.DB, c models.Commits) *gorm.DB {
	return db.Create(&c)
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