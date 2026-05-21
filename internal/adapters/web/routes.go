package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/ai"
	handlerCourse "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/enrollment"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/exercise"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/material"
	courselevel "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course_level"
	handlerNB "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/notebook"
	practicesheet "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/practice_sheet"
	studentprogress "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_progress"
	handlerTopic "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/topic"
	userprofile "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/user_profile"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/usecases"
)

func RegisterRoutes(app *gin.Engine, uc *usecases.Usecases) {
	api := app.Group("/api")
	api.Use(middlewares.AuthMiddleware())

	// Profile
	api.POST("/profile", userprofile.NewSyncHandler(uc.Profile.Sync))
	api.GET("/profile", userprofile.NewGetHandler(uc.Profile.Get))

	// Courses
	api.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	api.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))
	api.DELETE("/courses/:id", handlerCourse.NewDeleteHandler(uc.Course.Delete))

	// Enrollments
	api.POST("/courses/:id/enroll", enrollment.NewEnrollHandler(uc.Enrollment.Enroll))
	api.GET("/courses/:id/students", enrollment.NewListStudentsHandler(uc.Enrollment.ListStudents))

	// Materials
	api.POST("/courses/:id/materials", material.NewCreateHandler(uc.Material.Create))
	api.GET("/courses/:id/materials", material.NewListHandler(uc.Material.List))

	// Topics
	api.POST("/courses/:id/topics", handlerTopic.NewCreateHandler(uc.Topic.Create))
	api.GET("/courses/:id/topics", handlerTopic.NewListHandler(uc.Topic.List))

	// Exercises
	api.POST("/topics/:id/exercises", exercise.NewCreateHandler(uc.Exercise.Create))
	api.GET("/topics/:id/exercises", exercise.NewListHandler(uc.Exercise.List))
	api.PUT("/exercises/:id", exercise.NewUpdateHandler(uc.Exercise.Update))
	api.DELETE("/exercises/:id", exercise.NewDeleteHandler(uc.Exercise.Delete))

	// Practice Sheets
	api.POST("/courses/:id/practice-sheets", practicesheet.NewCreateHandler(uc.PracticeSheet.Create))
	api.GET("/courses/:id/practice-sheets", practicesheet.NewListHandler(uc.PracticeSheet.List))
	api.GET("/practice-sheets/:id", practicesheet.NewGetHandler(uc.PracticeSheet.Get))
	api.POST("/practice-sheets/:id/submit", practicesheet.NewSubmitHandler(uc.PracticeSheet.Submit))

	// Student Progress
	api.GET("/students/me/progress", studentprogress.NewGetMyProgressHandler(uc.Progress.GetMy))
	api.GET("/students/me/courses/:id/progress", studentprogress.NewGetCourseProgressHandler(uc.Progress.GetCourse))

	// AI Tutor
	api.POST("/ai/conversations", ai.NewCreateConversationHandler(uc.AI.CreateConversation))
	api.GET("/ai/conversations/:id/messages", ai.NewGetMessagesHandler(uc.AI.GetMessages))
	api.POST("/ai/help", ai.NewHelpHandler(uc.AI.Help))

	// Course levels
	api.GET("/courses/:id/levels", courselevel.NewGetHandler(uc.CourseLevel.Get))

	// Notebooks
	api.POST("/courses/:id/notebooks", handlerNB.NewCreateHandler(uc.Notebook.Create))
	api.GET("/courses/:id/notebooks", handlerNB.NewListHandler(uc.Notebook.List))
	api.GET("/notebooks/:id", handlerNB.NewGetHandler(uc.Notebook.Get))
	api.POST("/notebooks/:id/pages", handlerNB.NewAddPageHandler(uc.Notebook.AddPage))
	api.PUT("/notebook-pages/:id", handlerNB.NewUpdatePageHandler(uc.Notebook.UpdatePage))
	api.POST("/notebook-pages/:id/submit", handlerNB.NewSaveSubmissionHandler(uc.Notebook.SaveSubmission))
}
