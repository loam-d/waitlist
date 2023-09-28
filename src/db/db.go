package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"waitlist/src/models"
)

func StartDb() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("waitlist.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Student{})
	if err != nil {
		panic("failed to migrate database")
	}
	err = db.AutoMigrate(&models.Course{})
	if err != nil {
		panic("failed to migrate database")
	}
	err = db.AutoMigrate(&models.Enrollment{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db, err
}
