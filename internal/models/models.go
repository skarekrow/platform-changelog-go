package models

import (
	"time"
)

type Services struct {
	ID          int    `gorm:"primary_key;autoincremement"`
	Name        string `gorm:"not null"`
	DisplayName string `gorm:"not null;unique"`
	GHRepo      string
	GLRepo      string
	DeployFile  string
	Namespace   string
	Branch      string      `gorm:"default:'master'"`
	Timelines   []Timelines `gorm:"foreignkey:ServiceID"`
}

type timelineType string

const (
	commit timelineType = "commit"
	deploy timelineType = "deploy"
)

type Timelines struct {
	ID              int          `gorm:"primary_key;autoincrement" json:"id"`
	ServiceID       int          `gorm:"not null;foreign_key:services.id" json:"service_id"`
	Timestamp       time.Time    `gorm:"not null" json:"timestamp"`
	Type            timelineType `gorm:"not null" json:"type" sql:"type:timeline_type"`
	Repo            string       `gorm:"not null" json:"repo"`
	Ref             string       `json:"ref,omitempty"`
	Author          string       `json:"author,omitempty"`
	MergedBy        string       `json:"merged_by,omitempty"`
	Message         string       `json:"message,omitempty"`
	DeployNamespace string       `json:"namespace,omitempty"`
	Cluster         string       `json:"cluster,omitempty"`
	Image           string       `json:"image,omitempty"`
	TriggeredBy     string       `json:"triggered_by,omitempty"`
	Status          string       `json:"status,omitempty"`
}
