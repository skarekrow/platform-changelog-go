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
	Branch      string    `gorm:"default:'master'"`
	Commits     []Commits `gorm:"foreignkey:ID"`
	Deploys     []Deploys `gorm:"foreignkey:ID"`
}

type Commits struct {
	ID        int       `gorm:"primary_key;autoincrement"`
	ServiceID int       `gorm:"not null;foreign_key:services.id"`
	Ref       string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Repo      string    `gorm:"not null"`
	Author    string    `gorm:"not null"`
	MergedBy  string
	Message   string
}

type Deploys struct {
	ID        int       `gorm:"primary_key;autoincrement"`
	ServiceID int       `gorm:"not null;foreign_key:services.id"`
	Ref       string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Repo      string    `gorm:"not null"`
	Namespace string    `gorm:"not null"`
	Cluster   string    `gorm:"not null"`
	Image     string    `gorm:"not null"`
}

/**
 * A combination of the Commits and Deploys models
 */
type Timelines struct {
	ID        int       `gorm:"primary_key;autoincrement"`
	ServiceID int       `gorm:"not null;foreign_key:services.id"`
	CommitID  int       `gorm:"foreign_key:commits.id"`
	DeployID  int       `gorm:"foreign_key:deploys.id"`
	Timestamp time.Time `gorm:"not null"`
	Commit    Commits
	Deploy    Deploys
}
