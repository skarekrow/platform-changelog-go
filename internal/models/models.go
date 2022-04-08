package models

import (
	"time"
)

type Service struct {
	ID int `gorm:"primary_key"`
	Name string `gorm:"not null;unique"`
	DisplayName string
	GHRepo string
	GLRepo string
	DeployFile string
	Namespace string
	Branch string `gorm:"default:'master'"`
	Commits []Commit
	Deploys []Deploy
}

type Commit struct {
	ID int `gorm:"primary_key"`
	ServiceID int `gorm:"not null;foreign_key:service.id"`
	Repo string `gorm:"not null"`
	Ref string `gorm:"not null"`
	Title string `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	Author string `gorm:"not null"`
	Message string
}

type Deploy struct {
	ID int `gorm:"primary_key"`
	ServiceID int `gorm:"not null;foreign_key:service.id"`
	Ref string `gorm:"not null"`
	Namespace string `gorm:"not null"`
	Cluster string `gorm:"not null"`
	Image string `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}