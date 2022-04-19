package models

import (
	"time"
)

type Services struct {
	ID          int    `gorm:"primary_key"`
	DisplayName string `gorm:"not null;unique"`
	GHRepo      string
	GLRepo      string
	DeployFile  string
	Namespace   string
	Branch      string `gorm:"default:'master'"`
	Commits     Commits `gorm:"foreignkey:ID"`
	Deploys     Deploys `gorm:"foreignkey:ID"`
}

type Commits struct {
	ID        int       `gorm:"primary_key;autoincrement"`
	ServiceID int       `gorm:"not null;foreign_key:services.id"`
	Repo      string    `gorm:"not null"`
	Ref       string    `gorm:"not null"`
	Title     string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Author    string    `gorm:"not null"`
	MergedBy  string
	Message   string
}

type Deploys struct {
	ID        int       `gorm:"primary_key"`
	ServiceID int       `gorm:"not null;foreign_key:services.id"`
	Ref       string    `gorm:"not null"`
	Namespace string    `gorm:"not null"`
	Cluster   string    `gorm:"not null"`
	Image     string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}