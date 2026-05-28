package web

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/ai"
	handlerCourse "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course"
	courselevel "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course_level"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/enrollment"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/exercise"
	handlerGrade "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/grade"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/material"
	handlerNB "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/notebook"
	practicesheet "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/practice_sheet"
	studentprogress "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_progress"
	handlerSubject "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/subject"
	handlerAssignment "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/teacher_student_assignment"
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
	api.GET("/profile/:id", userprofile.NewGetByIDHandler(uc.Profile.Get))
	api.PUT("/profile/assistant-config", userprofile.NewUpdateAssistantConfigHandler(uc.Profile.UpdateAssistantConfig))
	adminOnly := api.Group("/")
	adminOnly.Use(middlewares.RequireRoles("admin", "superadmin"))
	adminOnly.PUT("/profile/:id/assistant-config", userprofile.NewUpdateAssistantConfigByIDHandler(uc.Profile.UpdateAssistantConfig))
	adminOnly.PUT("/profile/:id/academic-status", userprofile.NewUpdateAcademicStatusByIDHandler(uc.Profile.UpdateAcademicStatus))

	// Courses
	api.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	api.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))
	api.DELETE("/courses/:id", handlerCourse.NewDeleteHandler(uc.Course.Delete))

	// Grades
	api.POST("/grades", handlerGrade.NewCreateHandler(uc.Grade.Create))
	api.GET("/grades", handlerGrade.NewListHandler(uc.Grade.List))
	api.PUT("/grades/:id", handlerGrade.NewUpdateHandler(uc.Grade.Update))
	api.DELETE("/grades/:id", handlerGrade.NewDeleteHandler(uc.Grade.Delete))
	adminOnly.POST("/grades/:id/members", handlerGrade.NewAssignMemberHandler(uc.Grade.AssignMember))
	api.GET("/grades/:id/members", handlerGrade.NewListMembersHandler(uc.Grade.ListMembers))
	adminOnly.DELETE("/grades/:id/members/:userId", handlerGrade.NewRemoveMemberHandler(uc.Grade.RemoveMember))
	api.GET("/users/:userId/grades", handlerGrade.NewListUserGradesHandler(uc.Grade.ListUserGrades))

	// Subjects
	api.POST("/subjects", handlerSubject.NewCreateHandler(uc.Subject.Create))
	api.GET("/subjects", handlerSubject.NewListHandler(uc.Subject.List))
	api.PUT("/subjects/:id", handlerSubject.NewUpdateHandler(uc.Subject.Update))
	api.DELETE("/subjects/:id", handlerSubject.NewDeleteHandler(uc.Subject.Delete))

	// Teacher/student assignments
	adminOnly.POST("/teacher-student-assignments", handlerAssignment.NewAssignHandler(uc.Assignment.Assign))
	adminOnly.DELETE("/teacher-student-assignments/:teacherId/:studentId", handlerAssignment.NewUnassignHandler(uc.Assignment.Unassign))
	adminOnly.GET("/teachers/:teacherId/students", handlerAssignment.NewListStudentsHandler(uc.Assignment.ListStudents))
	adminOnly.GET("/students/:studentId/teachers", handlerAssignment.NewListTeachersHandler(uc.Assignment.ListTeachers))
	api.GET("/teachers/me/students", handlerAssignment.NewListMyStudentsHandler(uc.Assignment.ListStudents))

	// Enrollments
	api.POST("/courses/:id/enroll", enrollment.NewEnrollHandler(uc.Enrollment.Enroll))
	api.GET("/courses/:id/students", enrollment.NewListStudentsHandler(uc.Enrollment.ListStudents))

	// Materials
	api.POST("/courses/:id/materials", material.NewCreateHandler(uc.Material.Create))
	api.GET("/courses/:id/materials", material.NewListHandler(uc.Material.List))
	api.PUT("/materials/:id", material.NewUpdateHandler(uc.Material.Update))
	api.DELETE("/materials/:id", material.NewDeleteHandler(uc.Material.Delete))

	// Topics
	api.POST("/courses/:id/topics", handlerTopic.NewCreateHandler(uc.Topic.Create))
	api.GET("/courses/:id/topics", handlerTopic.NewListHandler(uc.Topic.List))
	api.PUT("/topics/:id", handlerTopic.NewUpdateHandler(uc.Topic.Update))
	api.DELETE("/topics/:id", handlerTopic.NewDeleteHandler(uc.Topic.Delete))

	// Exercises
	api.POST("/topics/:id/exercises", exercise.NewCreateHandler(uc.Exercise.Create))
	api.GET("/topics/:id/exercises", exercise.NewListHandler(uc.Exercise.List))
	api.PUT("/exercises/:id", exercise.NewUpdateHandler(uc.Exercise.Update))
	api.DELETE("/exercises/:id", exercise.NewDeleteHandler(uc.Exercise.Delete))

	// Practice Sheets
	api.POST("/courses/:id/practice-sheets", practicesheet.NewCreateHandler(uc.PracticeSheet.Create))
	api.GET("/courses/:id/practice-sheets", practicesheet.NewListHandler(uc.PracticeSheet.List))
	api.GET("/practice-sheets/:id", practicesheet.NewGetHandler(uc.PracticeSheet.Get))
	api.PUT("/practice-sheets/:id", practicesheet.NewUpdateHandler(uc.PracticeSheet.Update))
	api.DELETE("/practice-sheets/:id", practicesheet.NewDeleteHandler(uc.PracticeSheet.Delete))
	api.POST("/practice-sheets/:id/submit", practicesheet.NewSubmitHandler(uc.PracticeSheet.Submit))
	api.POST("/practice-sheets/:id/submit-async", practicesheet.NewSubmitAsyncHandler(uc.PracticeSheet.Submit))
	api.GET("/practice-sheets/submit-jobs/:jobId", practicesheet.NewGetSubmitJobHandler())

	// Student Progress (self-service)
	api.GET("/students/me/progress", studentprogress.NewGetMyProgressHandler(uc.Progress.GetMy))
	api.GET("/students/me/courses/:id/progress", studentprogress.NewGetCourseProgressHandler(uc.Progress.GetCourse))

	// Teacher view of student progress
	api.GET("/teachers/me/students/:studentId/progress", studentprogress.NewGetStudentProgressHandler(uc.Progress.GetStudentProgress))
	api.GET("/teachers/me/students/:studentId/courses/:courseId/progress", studentprogress.NewGetStudentCourseProgressHandler(uc.Progress.GetStudentCourseProgress))
	api.GET("/teachers/me/students/:studentId/attempts", studentprogress.NewGetStudentAttemptsHandler(uc.Progress.GetStudentAttempts))

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
	api.PUT("/notebooks/:id", handlerNB.NewUpdateHandler(uc.Notebook.Update))
	api.DELETE("/notebooks/:id", handlerNB.NewDeleteHandler(uc.Notebook.Delete))
	api.POST("/notebooks/:id/pages", handlerNB.NewAddPageHandler(uc.Notebook.AddPage))
	api.PUT("/notebook-pages/:id", handlerNB.NewUpdatePageHandler(uc.Notebook.UpdatePage))
	api.POST("/notebook-pages/:id/submit", handlerNB.NewSaveSubmissionHandler(uc.Notebook.SaveSubmission))
	api.POST("/notebook-pages/:id/submit-async", handlerNB.NewSaveSubmissionAsyncHandler(uc.Notebook.SaveSubmission))
	api.GET("/notebook-pages/submit-jobs/:jobId", handlerNB.NewGetSubmitJobHandler())
}
