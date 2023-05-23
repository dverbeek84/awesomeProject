package model

import (
	"gorm.io/gorm"
)

type Deployment struct {
	gorm.Model
	State         string `json:"state"`
	ApplicationID int
	Application   ApplicationTemplate
}

type ApplicationTemplate struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	ExecCommand string `gorm:"not null" json:"exec_command"`
}
