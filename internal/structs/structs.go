package structs

import "github.com/redhatinsights/platform-changelog-go/internal/models"

type Query struct {
	Offset int
	Limit  int
}
type ServicesList struct {
	Count int64          `json:"count"`
	Data  []ServicesData `json:"data"`
}

type ExpandedServicesList struct {
	Count int64                  `json:"count"`
	Data  []ExpandedServicesData `json:"data"`
}

type TimelinesList struct {
	Count int64              `json:"count"`
	Data  []models.Timelines `json:"data"`
}

type ServicesData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	GHRepo      string `json:"gh_repo"`
	GLRepo      string `json:"gl_repo"`
	DeployFile  string `json:"deploy_file"`
	Namespace   string `json:"namespace"`
	Branch      string `json:"branch"`
}

type ExpandedServicesData struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	GHRepo      string           `json:"gh_repo"`
	GLRepo      string           `json:"gl_repo"`
	DeployFile  string           `json:"deploy_file"`
	Namespace   string           `json:"namespace"`
	Branch      string           `json:"branch"`
	Commit      models.Timelines `json:"latest_commit" gorm:"foreignkey:ID"`
	Deploy      models.Timelines `json:"latest_deploy" gorm:"foreignkey:ID"`
}
