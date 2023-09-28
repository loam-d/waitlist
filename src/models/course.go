package models

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);unique_index"`
}

func CreateCourse(db *gorm.DB, courseName string) (course *Course, err error) {
	course = &Course{Name: courseName}

	result := db.Create(course)
	if result.Error != nil {
		return nil, result.Error
	}

	return course, nil
}
