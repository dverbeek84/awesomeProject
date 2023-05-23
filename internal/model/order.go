package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	State         string `json:"State"`
	ApplicationID int    `json:"ApplicationID" binding:"required"`
}
