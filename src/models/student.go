package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email string `gorm:"type:varchar(100);unique_index"`
}

func CreateStudent(db *gorm.DB, email string) (student *Student, err error) {
	student = &Student{Email: email}

	result := db.Create(student)
	if result.Error != nil {
		return nil, result.Error
	}

	return student, nil
}
