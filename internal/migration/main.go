package main

import (
	"github.com/redhatinsights/platform-changelog-go/internal/config"
	"github.com/redhatinsights/platform-changelog-go/internal/db"
	"github.com/redhatinsights/platform-changelog-go/internal/logging"
	"github.com/redhatinsights/platform-changelog-go/internal/models"
)

func main() {
	logging.InitLogger()

	cfg := config.Get()

	db.DbConnect(cfg)

	db.DB.AutoMigrate(
		&models.Service{},
		&models.Commit{},
		&models.Deploy{},
	)

	logging.Log.Info("DB Migration Complete")
}