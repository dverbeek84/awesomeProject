package database

import (
	"database/sql"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(name string) (*sql.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(name), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return DB.DB()
}
