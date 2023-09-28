package models

import (
	"gorm.io/gorm"
	"time"
)

type enrollmentStatus string

const (
	WAITING   enrollmentStatus = "WAITING"
	ENROLLED  enrollmentStatus = "ENROLLED"
	OFFERED   enrollmentStatus = "OFFERED"
	REJECTED  enrollmentStatus = "REJECTED"
	COMPLETED enrollmentStatus = "COMPLETED"
)

const TIME_OUT_DAYS = 2

type Enrollment struct {
	gorm.Model
	StudentID   uint             `gorm:"primaryKey"`
	CourseID    uint             `gorm:"primaryKey"`
	Status      enrollmentStatus `gorm:"type:enrollment_status"`
	QueuedAt    time.Time
	OfferedAt   time.Time
	CompletedAt time.Time
	Student     Student `gorm:"foreignKey:StudentID"`
	Course      Course  `gorm:"foreignKey:CourseID"`
}

func AddStudentToCourseWaitlist(db *gorm.DB, student *Student, course *Course) (*Enrollment, error) {
	enrollment := &Enrollment{
		StudentID: student.ID,
		CourseID:  course.ID,
		Status:    WAITING,
		QueuedAt:  time.Now(),
	}

	result := db.Create(enrollment)
	if result.Error != nil {
		return nil, result.Error
	}

	return enrollment, nil
}

func OfferCourseToStudents(db *gorm.DB, course *Course, numSlots int) ([]Enrollment, error) {
	var enrollmentsToOffer []Enrollment
	result := db.Order("queued_at ASC").Limit(numSlots).Where("status = ?", WAITING).Find(&enrollmentsToOffer)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, enrollment := range enrollmentsToOffer {
		result := db.Model(&enrollment).Updates(Enrollment{Status: OFFERED, OfferedAt: time.Now()})
		if result.Error != nil {
			return nil, result.Error
		}
		// TODO: Send an email for student to accept offer
	}

	return enrollmentsToOffer, nil
}

func AcceptOffer(db *gorm.DB, enrollment *Enrollment) (*Enrollment, error) {
	result := db.Model(&enrollment).Update("Status", ENROLLED)
	if result.Error != nil {
		return nil, result.Error
	}
	return enrollment, nil
}

func RejectOffer(db *gorm.DB, enrollment *Enrollment) (*Enrollment, error) {
	course := enrollment.Course
	newEnrollment, err := OfferCourseToStudents(db, &course, 1)
	if err != nil {
		return nil, err
	}
	result := db.Model(&enrollment).Update("Status", REJECTED)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(newEnrollment) > 0 {
		return &newEnrollment[0], nil
	} else {
		return nil, nil
	}
}

func TimeOutOffers(db *gorm.DB, requeue bool) ([]Enrollment, error) {
	var enrollmentsToTimeOut []Enrollment

	daysAgo := time.Now().Add(-TIME_OUT_DAYS * 24 * time.Hour)
	result := db.Where("status = ? AND offered_at <= ?", OFFERED, daysAgo).Find(&enrollmentsToTimeOut)
	if result.Error != nil {
		return nil, result.Error
	}

	var newEnrollments []Enrollment
	for _, enrollment := range enrollmentsToTimeOut {
		var result *gorm.DB
		if requeue {
			result = db.Model(&enrollment).Updates(Enrollment{Status: WAITING, QueuedAt: time.Now()})
		} else {
			result = db.Model(&enrollment).Update("Status", REJECTED)
		}
		if result.Error != nil {
			return nil, result.Error
		}

		newEnrollment, err := OfferCourseToStudents(db, &enrollment.Course, 1)
		if err != nil {
			return nil, result.Error
		}
		newEnrollments = append(newEnrollments, newEnrollment[0])
	}

	for _, newEnroll := range newEnrollments {
		enrollmentsToTimeOut = append(enrollmentsToTimeOut, newEnroll)
	}
	return enrollmentsToTimeOut, nil
}

func CompleteCourse(db *gorm.DB, course *Course) ([]Enrollment, error) {
	var enrollmentsToComplete []Enrollment
	result := db.Where("course_id = ? AND status = ?", course.ID, ENROLLED).Find(&enrollmentsToComplete)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, enrollment := range enrollmentsToComplete {
		result := db.Model(&enrollment).Updates(Enrollment{Status: COMPLETED, CompletedAt: time.Now()})
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return enrollmentsToComplete, nil
}
