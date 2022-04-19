package main

import (
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/models"

	"gorm.io/gorm"
)

func main() {
	logging.InitLogger()

	cfg := config.Get()

	db.DbConnect(cfg)

	db.DB.AutoMigrate(
		&models.Services{},
		&models.Commits{},
		&models.Deploys{},
	)

	logging.Log.Info("DB Migration Complete")

	reconcileServices(db.DB, cfg)
}

func reconcileServices(g *gorm.DB, cfg *config.Config) {
	for _, s := range cfg.Services {
		res, _ := db.GetServiceByName(g, s.DisplayName)
		if res.RowsAffected == 0 {
			_, service := db.CreateServiceTableEntry(g, s)
			logging.Log.Info("Created service: ", service)
		} else {
			logging.Log.Info("Service already exists: ", s.DisplayName)
		}
	}
}
