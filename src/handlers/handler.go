package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"waitlist/src/models"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) getEnrollment(courseName string, studentEmail string) (*models.Enrollment, error) {
	var course models.Course
	result := h.DB.Where("name = ?", courseName).First(&course)
	if result.Error != nil {
		return nil, result.Error
	}
	var student models.Student
	result = h.DB.Where("email = ?", studentEmail).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	var enrollment models.Enrollment
	result = h.DB.Where("student_id = ? AND course_id = ?", student.ID, course.ID).First(&enrollment)
	if result.Error != nil {
		return nil, result.Error
	}

	return &enrollment, nil
}

func (h *Handler) CreateCourse(c echo.Context) error {
	courseName := c.Param("name")
	_, err := models.CreateCourse(h.DB, courseName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "course creation failed")
	}
	return c.String(http.StatusCreated, fmt.Sprintf("Course %s added", courseName))
}

func (h *Handler) CreateStudent(c echo.Context) error {
	studentEmail := c.Param("email")
	_, err := models.CreateStudent(h.DB, studentEmail)
	if err != nil {
		return c.String(http.StatusInternalServerError, "student creation failed")
	}
	return c.String(http.StatusCreated, fmt.Sprintf("Student %s added", studentEmail))
}

func (h *Handler) AddStudentToCourseWaitlist(c echo.Context) error {
	courseName := c.Param("name")
	var course models.Course
	result := h.DB.Where("name = ?", courseName).First(&course)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "adding student to waitlist failed")
	}
	studentEmail := c.Param("email")
	var student models.Student
	result = h.DB.Where("email = ?", studentEmail).First(&student)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "adding student to waitlist failed")
	}
	_, err := models.AddStudentToCourseWaitlist(h.DB, &student, &course)
	if err != nil {
		return c.String(http.StatusInternalServerError, "adding student to waitlist failed")
	}
	return c.String(http.StatusCreated, fmt.Sprintf("student %s added to course %s waitlist", studentEmail, courseName))
}

func (h *Handler) OfferCourseToStudents(c echo.Context) error {
	courseName := c.Param("name")
	var course models.Course
	result := h.DB.Where("name = ?", courseName).First(&course)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "offering course failed")
	}
	numSlots, err := strconv.Atoi(c.Param("slots"))
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "offering course failed")
	}
	_, err = models.OfferCourseToStudents(h.DB, &course, numSlots)
	if err != nil {
		return c.String(http.StatusInternalServerError, "offering course failed")
	}
	return c.String(http.StatusOK, fmt.Sprintf("students offered course %s", courseName))
}

func (h *Handler) AcceptOffer(c echo.Context) error {
	courseName := c.Param("name")
	studentEmail := c.Param("email")
	enrollment, err := h.getEnrollment(courseName, studentEmail)
	if err != nil {
		return c.String(http.StatusInternalServerError, "accepting offer failed")
	}

	_, err = models.AcceptOffer(h.DB, enrollment)
	if err != nil {
		return c.String(http.StatusInternalServerError, "accepting offer failed")
	}
	return c.String(http.StatusOK, fmt.Sprintf("student %s accepted course %s offer", studentEmail, courseName))
}

func (h *Handler) RejectOffer(c echo.Context) error {
	courseName := c.Param("name")
	studentEmail := c.Param("email")
	enrollment, err := h.getEnrollment(courseName, studentEmail)
	if err != nil {
		return c.String(http.StatusInternalServerError, "rejecting offer failed")
	}

	_, err = models.RejectOffer(h.DB, enrollment)
	if err != nil {
		return c.String(http.StatusInternalServerError, "rejecting offer failed")
	}
	return c.String(http.StatusOK, fmt.Sprintf("student %s rejected course %s offer", studentEmail, courseName))
}

func (h *Handler) TimeOutOffers(c echo.Context) error {
	enrollments, err := models.TimeOutOffers(h.DB, true)
	if err != nil {
		return c.String(http.StatusInternalServerError, "timeout course failed")
	}
	return c.String(http.StatusOK, fmt.Sprintf("timed out %d students", len(enrollments)))
}

func (h *Handler) CompleteCourse(c echo.Context) error {
	courseName := c.Param("name")
	var course models.Course
	result := h.DB.Where("name = ?", courseName).First(&course)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "course completion failed")
	}

	enrollments, err := models.CompleteCourse(h.DB, &course)
	if err != nil {
		return c.String(http.StatusInternalServerError, "course completion failed")
	}
	return c.String(http.StatusOK, fmt.Sprintf("course %s had %d students complete it", courseName, len(enrollments)))
}
