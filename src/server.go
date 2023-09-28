package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"waitlist/src/db"
	"waitlist/src/handlers"
)

func main() {
	e := echo.New()

	database, err := db.StartDb()
	if err != nil {
		panic("DB did not start correctly")
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "TODO")
	})

	handler := &handlers.Handler{DB: database}
	e.POST("/course/:name", handler.CreateCourse)
	e.POST("/student/:email", handler.CreateStudent)
	e.POST("/course/:name/:email", handler.AddStudentToCourseWaitlist)
	e.PUT("/course/offer/:name/:slots", handler.OfferCourseToStudents)
	e.PUT("/student/accept/:email/:name", handler.AcceptOffer)
	e.PUT("/student/reject/:email/:name", handler.RejectOffer)
	e.PUT("/course/timeout/:name", handler.TimeOutOffers)
	e.PUT("/course/complete/:name", handler.CompleteCourse)
	e.GET("/course", handler.GetAllCourses)
	e.GET("/course/:name", handler.GetCourse)
	e.GET("/student", handler.GetAllStudents)
	e.GET("/student/:email", handler.GetStudent)
	e.GET("/enrollments", handler.GetAllEnrollments)

	e.Logger.Fatal(e.Start(":1323"))
}
