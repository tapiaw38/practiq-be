package web

import (
	"github.com/gin-gonic/gin"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/ai"
	handlerReview "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/attempt_review"
	handlerCourse "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course"
	courselevel "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course_level"
	handlerCP "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course_progress"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/enrollment"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/exercise"
	handlerGrade "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/grade"
	handlerLS "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/learning_strategy"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/material"
	handlerNB "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/notebook"
	handlerNotification "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/notification"
	practicesheet "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/practice_sheet"
	studentprogress "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_progress"
	studentreport "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_report"
	handlerSubject "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/subject"
	handlerAssignment "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/teacher_student_assignment"
	handlerTopic "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/topic"
	handlerUpload "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/upload"
	userprofile "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/user_profile"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/usecases"
)

func RegisterRoutes(app *gin.Engine, uc *usecases.Usecases, submitJobRepo submitjob.Repository) {
	api := app.Group("/api")
	api.Use(middlewares.AuthMiddleware())

	// Profile
	api.POST("/profile", userprofile.NewSyncHandler(uc.Profile.Sync))
	api.GET("/profile", userprofile.NewGetHandler(uc.Profile.Get))
	api.GET("/profile/:id", userprofile.NewGetByIDHandler(uc.Profile.Get))
	api.PUT("/profile/assistant-config", userprofile.NewUpdateAssistantConfigHandler(uc.Profile.UpdateAssistantConfig))
	// admin es el rol del profesor y superadmin el del administrador; "teacher"
	// no existe como rol en auth-api-be, así que exigirlo nunca coincidía y el
	// grupo quedaba abierto a cualquiera de los otros dos.
	adminOnly := api.Group("/")
	adminOnly.Use(middlewares.RequireRoles(middlewares.RoleSuperAdmin))
	teacherOnly := api.Group("/")
	teacherOnly.Use(middlewares.RequireRoles(middlewares.RoleTeacher, middlewares.RoleSuperAdmin))
	adminOnly.PUT("/profile/:id/assistant-config", userprofile.NewUpdateAssistantConfigByIDHandler(uc.Profile.UpdateAssistantConfig))
	adminOnly.PUT("/profile/:id/academic-status", userprofile.NewUpdateAcademicStatusByIDHandler(uc.Profile.UpdateAcademicStatus))

	// Courses. Los use cases validan que el curso sea del profesor; el grupo
	// evita además que un alumno llegue siquiera a intentarlo.
	teacherOnly.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	teacherOnly.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))
	teacherOnly.DELETE("/courses/:id", handlerCourse.NewDeleteHandler(uc.Course.Delete))

	// Grades. La estructura académica es institucional: la escribe el
	// administrador. El listado queda abierto porque el alumno arma con él su
	// propia navegación.
	adminOnly.POST("/grades", handlerGrade.NewCreateHandler(uc.Grade.Create))
	api.GET("/grades", handlerGrade.NewListHandler(uc.Grade.List))
	adminOnly.PUT("/grades/:id", handlerGrade.NewUpdateHandler(uc.Grade.Update))
	adminOnly.DELETE("/grades/:id", handlerGrade.NewDeleteHandler(uc.Grade.Delete))
	adminOnly.POST("/grades/:id/members", handlerGrade.NewAssignMemberHandler(uc.Grade.AssignMember))
	teacherOnly.GET("/grades/:id/members", handlerGrade.NewListMembersHandler(uc.Grade.ListMembers))
	adminOnly.DELETE("/grades/:id/members/:userId", handlerGrade.NewRemoveMemberHandler(uc.Grade.RemoveMember))
	api.GET("/users/:userId/grades", handlerGrade.NewListUserGradesHandler(uc.Grade.ListUserGrades))

	// Subjects. Mismo criterio que grades.
	adminOnly.POST("/subjects", handlerSubject.NewCreateHandler(uc.Subject.Create))
	api.GET("/subjects", handlerSubject.NewListHandler(uc.Subject.List))
	adminOnly.PUT("/subjects/:id", handlerSubject.NewUpdateHandler(uc.Subject.Update))
	adminOnly.DELETE("/subjects/:id", handlerSubject.NewDeleteHandler(uc.Subject.Delete))

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
	teacherOnly.POST("/courses/:id/materials", material.NewCreateHandler(uc.Material.Create))
	api.GET("/courses/:id/materials", material.NewListHandler(uc.Material.List))
	api.GET("/materials/:id", material.NewGetHandler(uc.Material.Get))
	teacherOnly.PUT("/materials/:id", material.NewUpdateHandler(uc.Material.Update))
	teacherOnly.DELETE("/materials/:id", material.NewDeleteHandler(uc.Material.Delete))

	// Topics
	teacherOnly.POST("/courses/:id/topics", handlerTopic.NewCreateHandler(uc.Topic.Create))
	api.GET("/courses/:id/topics", handlerTopic.NewListHandler(uc.Topic.List))
	teacherOnly.PUT("/topics/:id", handlerTopic.NewUpdateHandler(uc.Topic.Update))
	teacherOnly.DELETE("/topics/:id", handlerTopic.NewDeleteHandler(uc.Topic.Delete))

	// Exercises
	teacherOnly.POST("/topics/:id/exercises", exercise.NewCreateHandler(uc.Exercise.Create))
	api.GET("/topics/:id/exercises", exercise.NewListHandler(uc.Exercise.List))
	api.GET("/exercises/:id/statement-image", exercise.NewStatementImageHandler(uc.Exercise.StatementImage))
	teacherOnly.PUT("/exercises/:id", exercise.NewUpdateHandler(uc.Exercise.Update))
	teacherOnly.DELETE("/exercises/:id", exercise.NewDeleteHandler(uc.Exercise.Delete))

	// Practice Sheets
	teacherOnly.POST("/courses/:id/practice-sheets", practicesheet.NewCreateHandler(uc.PracticeSheet.Create))
	api.GET("/courses/:id/practice-sheets", practicesheet.NewListHandler(uc.PracticeSheet.List))
	api.GET("/practice-sheets/:id", practicesheet.NewGetHandler(uc.PracticeSheet.Get))
	api.GET("/practice-sheets/:id/exercises/:exerciseId/assistant-media", practicesheet.NewGetAssistantMediaHandler(uc.PracticeSheet.GetAssistantMedia))
	teacherOnly.PUT("/practice-sheets/:id", practicesheet.NewUpdateHandler(uc.PracticeSheet.Update))
	teacherOnly.DELETE("/practice-sheets/:id", practicesheet.NewDeleteHandler(uc.PracticeSheet.Delete))
	api.POST("/practice-sheets/:id/submit", practicesheet.NewSubmitHandler(uc.PracticeSheet.Submit))
	api.POST("/practice-sheets/:id/submit-async", practicesheet.NewSubmitAsyncHandler(uc.PracticeSheet.Submit, submitJobRepo))
	api.GET("/practice-sheets/submit-jobs/:jobId", practicesheet.NewGetSubmitJobHandler(submitJobRepo))

	// Student Progress (self-service)
	api.GET("/students/me/progress", studentprogress.NewGetMyProgressHandler(uc.Progress.GetMy))
	api.GET("/students/me/dashboard", studentprogress.NewDashboardHandler(uc.Progress.Dashboard))
	api.GET("/students/me/courses/:id/progress", studentprogress.NewGetCourseProgressHandler(uc.Progress.GetCourse))

	// Teacher view of student progress
	api.GET("/teachers/me/students/:studentId/progress", studentprogress.NewGetStudentProgressHandler(uc.Progress.GetStudentProgress))
	api.GET("/teachers/me/students/:studentId/courses/:courseId/progress", studentprogress.NewGetStudentCourseProgressHandler(uc.Progress.GetStudentCourseProgress))
	api.GET("/teachers/me/students/:studentId/attempts", studentprogress.NewGetStudentAttemptsHandler(uc.Progress.GetStudentAttempts))
	teacherOnly.GET("/teachers/me/students/:studentId/report.pdf", studentreport.NewGeneratePDFHandler(uc.Report.GeneratePDF))

	// AI Tutor
	api.POST("/ai/conversations", ai.NewCreateConversationHandler(uc.AI.CreateConversation))
	api.GET("/ai/conversations/:id/messages", ai.NewGetMessagesHandler(uc.AI.GetMessages))
	api.POST("/ai/help", ai.NewHelpHandler(uc.AI.Help))
	api.POST("/ai/copilot", ai.NewCopilotHandler(uc.AI.Help))
	api.POST("/ai/copilot/stream", ai.NewCopilotStreamHandler(uc.AI.Help))
	api.POST("/ai/curiosities", ai.NewGenerateCuriositiesHandler(uc.AI.GenerateCuriosities))
	api.GET("/assistant-proxy/conversation/user", ai.NewProxyListConversationsHandler(uc.AI.Proxy))
	api.GET("/assistant-proxy/conversation/:id", ai.NewProxyGetConversationHandler(uc.AI.Proxy))
	api.POST("/assistant-proxy/conversation/", ai.NewProxyCreateConversationHandler(uc.AI.Proxy))
	api.POST("/assistant-proxy/conversation/:id/message", ai.NewProxySendMessageHandler(uc.AI.Proxy))
	api.POST("/assistant-proxy/conversation/:id/message/text", ai.NewProxySendTextMessageHandler(uc.AI.Proxy))

	// Course levels
	api.GET("/courses/:id/levels", courselevel.NewGetHandler(uc.CourseLevel.Get))

	// Notebooks
	teacherOnly.POST("/courses/:id/notebooks", handlerNB.NewCreateHandler(uc.Notebook.Create))
	api.GET("/courses/:id/notebooks", handlerNB.NewListHandler(uc.Notebook.List))
	api.GET("/notebooks/:id", handlerNB.NewGetHandler(uc.Notebook.Get))
	teacherOnly.PUT("/notebooks/:id", handlerNB.NewUpdateHandler(uc.Notebook.Update))
	teacherOnly.DELETE("/notebooks/:id", handlerNB.NewDeleteHandler(uc.Notebook.Delete))
	teacherOnly.POST("/notebooks/:id/pages", handlerNB.NewAddPageHandler(uc.Notebook.AddPage))
	teacherOnly.PUT("/notebook-pages/:id", handlerNB.NewUpdatePageHandler(uc.Notebook.UpdatePage))
	api.POST("/notebook-pages/:id/submit", handlerNB.NewSaveSubmissionHandler(uc.Notebook.SaveSubmission))
	api.POST("/notebook-pages/:id/submit-async", handlerNB.NewSaveSubmissionAsyncHandler(uc.Notebook.SaveSubmission, submitJobRepo))
	api.GET("/notebook-pages/submit-jobs/:jobId", handlerNB.NewGetSubmitJobHandler(submitJobRepo))
	teacherOnly.GET("/notebook-submissions", handlerNB.NewListSubmissionsHandler(uc.Notebook.ListSubmissions))
	teacherOnly.POST("/notebook-submissions/:id/review", handlerNB.NewReviewSubmissionHandler(uc.Notebook.ReviewSubmission))
	teacherOnly.PUT("/notebook-submissions/:id/teacher-review", handlerNB.NewTeacherReviewSubmissionHandler(uc.Notebook.TeacherReview))

	// Attachment answers the assistant could not grade
	teacherOnly.GET("/attempt-reviews", handlerReview.NewListHandler(uc.AttemptReview.List))
	teacherOnly.POST("/attempt-reviews/:id", handlerReview.NewReviewHandler(uc.AttemptReview.Review))
	teacherOnly.GET("/attempt-reviews/:id/statement-image", handlerReview.NewStatementImageHandler(uc.AttemptReview.StatementImage))

	// File uploads (attachment answers, teacher materials)
	api.POST("/uploads", handlerUpload.NewHandler(uc.Upload.Upload))

	// Notifications
	api.GET("/notifications", handlerNotification.NewListHandler(uc.Notification.List))
	api.POST("/notifications/:id/read", handlerNotification.NewMarkReadHandler(uc.Notification.MarkRead))
	api.POST("/notifications/read-all", handlerNotification.NewMarkAllReadHandler(uc.Notification.MarkAllRead))
	api.DELETE("/notifications/:id", handlerNotification.NewDeleteHandler(uc.Notification.Delete))

	// Learning Strategies
	api.GET("/learning-strategies", handlerLS.NewListHandler(uc.LearningStrategy.List))
	api.GET("/learning-strategies/:id", handlerLS.NewGetHandler(uc.LearningStrategy.Get))
	adminOnly.POST("/learning-strategies", handlerLS.NewCreateHandler(uc.LearningStrategy.Create))
	adminOnly.PUT("/learning-strategies/:id", handlerLS.NewUpdateHandler(uc.LearningStrategy.Update))
	adminOnly.DELETE("/learning-strategies/:id", handlerLS.NewDeleteHandler(uc.LearningStrategy.Delete))

	// Course Learning Strategies
	api.GET("/courses/:id/strategies", handlerLS.NewListByCourseHandler(uc.LearningStrategy.ListByCourse))
	teacherOnly.POST("/courses/:id/strategies", handlerLS.NewAssignToCourseHandler(uc.LearningStrategy.AssignToCourse))
	teacherOnly.DELETE("/course-learning-strategies/:id", handlerLS.NewUnassignFromCourseHandler(uc.LearningStrategy.UnassignFromCourse))

	// Course Progress
	teacherOnly.GET("/students/:studentId/courses/:courseId/progress", handlerCP.NewGetForStudentHandler(uc.CourseProgress.GetForStudent))
	teacherOnly.GET("/students/:studentId/progress", handlerCP.NewListForStudentHandler(uc.CourseProgress.ListForStudent))
}
