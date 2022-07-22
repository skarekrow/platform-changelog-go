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
	Timelines   []Timelines `gorm:"foreignkey:ID"`
}

// type timelineType string

// const (
// 	commit timelineType = "commit"
// 	deploy timelineType = "deploy"
// )

type Timelines struct {
	ID        int       `gorm:"primary_key;autoincrement"`
	ServiceID int       `gorm:"not null;foreign_key:services.id"`
	Timestamp time.Time `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Ref       string    `gorm:"not null"`
	Repo      string    `gorm:"not null"`
	Author    string
	MergedBy  string
	Message   string
	Namespace string
	Cluster   string
	Image     string
}
