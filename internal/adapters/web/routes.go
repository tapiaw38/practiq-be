package web

import (
	"github.com/gin-gonic/gin"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	userprofileRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/user_profile"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/ai"
	handlerReview "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/attempt_review"
	handlerCourse "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course"
	courselevel "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course_level"
	handlerCP "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/course_progress"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/enrollment"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/exercise"
	gilliesettings "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/gillie_settings"
	handlerGrade "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/grade"
	handlerLS "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/learning_strategy"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/material"
	handlerNB "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/notebook"
	handlerNotification "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/notification"
	practicesheet "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/practice_sheet"
	handlerSchool "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/school"
	sitecontact "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/site_contact"
	handlerInvitation "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_invitation"
	studentprogress "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_progress"
	studentreport "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/student_report"
	handlerSubject "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/subject"
	subscription "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/subscription"
	handlerAssignment "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/teacher_student_assignment"
	handlerTopic "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/topic"
	handlerUpload "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/upload"
	userprofile "github.com/tapiaw38/practiq-be/internal/adapters/web/handlers/user_profile"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/platform/config"
	"github.com/tapiaw38/practiq-be/internal/platform/revocation"
	"github.com/tapiaw38/practiq-be/internal/usecases"
	ucSubscription "github.com/tapiaw38/practiq-be/internal/usecases/subscription"
)

func RegisterRoutes(app *gin.Engine, uc *usecases.Usecases, submitJobRepo submitjob.Repository, userProfiles userprofileRepo.Repository, revoked *revocation.Checker) {

	public := app.Group("/api/public")
	public.GET("/subscription-plans", subscription.NewPublicListPlansHandler(uc.Subscription.ListPlans))
	public.GET("/site-contact", sitecontact.NewGetHandler(uc.SiteContact.Get))

	api := app.Group("/api")
	api.Use(middlewares.AuthMiddleware(revoked))
	api.Use(middlewares.LoadProfileType(userProfiles))

	api.POST("/profile", userprofile.NewSyncHandler(uc.Profile.Sync))
	api.GET("/profile", userprofile.NewGetHandler(uc.Profile.Get))
	api.GET("/profile/:id", userprofile.NewGetByIDHandler(uc.Profile.Get))
	api.PUT("/profile/ui-theme", userprofile.NewUpdateUIThemeHandler(uc.Profile.UpdateUITheme))
	api.PUT("/profile/avatar", userprofile.NewUpdateAvatarSeedHandler(uc.Profile.UpdateAvatarSeed))
	adminOnly := api.Group("/")
	adminOnly.Use(middlewares.RequireRoles(middlewares.RoleSuperAdmin))
	adminOnly.GET("/site-contact", sitecontact.NewGetHandler(uc.SiteContact.Get))
	adminOnly.PUT("/site-contact", sitecontact.NewUpdateHandler(uc.SiteContact.Update))

	adminOnly.GET("/gillie-settings", gilliesettings.NewGetHandler(uc.GillieSettings.Get))
	adminOnly.PUT("/gillie-settings", gilliesettings.NewUpdateHandler(uc.GillieSettings.Update))
	teacherOnly := api.Group("/")
	teacherOnly.Use(middlewares.RequireTeacher())

	teacherOnly.PUT("/profile/:id/ui-theme", userprofile.NewUpdateUIThemeByIDHandler(uc.Profile.UpdateUITheme))
	teacherOnly.GET("/profile/find-by-email", userprofile.NewFindByEmailHandler(uc.Profile.FindByEmail))
	adminOnly.PUT("/profile/:id/academic-status", userprofile.NewUpdateAcademicStatusByIDHandler(uc.Profile.UpdateAcademicStatus))
	adminOnly.PUT("/profile/:id/type", userprofile.NewUpdateProfileTypeByIDHandler(uc.Profile.UpdateProfileType))

	teacherOnly.POST("/courses", handlerCourse.NewCreateHandler(uc.Course.Create))
	api.GET("/courses", handlerCourse.NewListHandler(uc.Course.List))
	api.GET("/courses/:id", handlerCourse.NewGetHandler(uc.Course.Get))
	teacherOnly.PUT("/courses/:id", handlerCourse.NewUpdateHandler(uc.Course.Update))
	teacherOnly.PATCH("/courses/:id/status", handlerCourse.NewSetStatusHandler(uc.Course.SetStatus))
	teacherOnly.DELETE("/courses/:id", handlerCourse.NewDeleteHandler(uc.Course.Delete))

	teacherOnly.POST("/grades", handlerGrade.NewCreateHandler(uc.Grade.Create))
	api.GET("/grades", handlerGrade.NewListHandler(uc.Grade.List))
	teacherOnly.PUT("/grades/:id", handlerGrade.NewUpdateHandler(uc.Grade.Update))
	teacherOnly.DELETE("/grades/:id", handlerGrade.NewDeleteHandler(uc.Grade.Delete))
	teacherOnly.POST("/grades/:id/members", handlerGrade.NewAssignMemberHandler(uc.Grade.AssignMember))
	teacherOnly.GET("/grades/:id/members", handlerGrade.NewListMembersHandler(uc.Grade.ListMembers))
	teacherOnly.DELETE("/grades/:id/members/:userId", handlerGrade.NewRemoveMemberHandler(uc.Grade.RemoveMember))
	api.GET("/users/:userId/grades", handlerGrade.NewListUserGradesHandler(uc.Grade.ListUserGrades))

	api.POST("/grades/batch-by-users", handlerGrade.NewListGradesByUsersHandler(uc.Grade.ListGradesByUsers))

	teacherOnly.POST("/subjects", handlerSubject.NewCreateHandler(uc.Subject.Create))
	api.GET("/subjects", handlerSubject.NewListHandler(uc.Subject.List))
	teacherOnly.PUT("/subjects/:id", handlerSubject.NewUpdateHandler(uc.Subject.Update))
	teacherOnly.DELETE("/subjects/:id", handlerSubject.NewDeleteHandler(uc.Subject.Delete))

	teacherOnly.POST("/teacher-student-assignments", handlerAssignment.NewAssignHandler(uc.Assignment.Assign))
	teacherOnly.DELETE("/teacher-student-assignments/:teacherId/:studentId", handlerAssignment.NewUnassignHandler(uc.Assignment.Unassign))
	teacherOnly.GET("/teachers/:teacherId/students", handlerAssignment.NewListStudentsHandler(uc.Assignment.ListStudents))
	teacherOnly.GET("/students/:studentId/teachers", handlerAssignment.NewListTeachersHandler(uc.Assignment.ListTeachers))
	api.GET("/teachers/me/students", handlerAssignment.NewListMyStudentsHandler(uc.Assignment.ListStudents))

	teacherOnly.POST("/invitations", handlerInvitation.NewCreateHandler(uc.Invitation.Create))
	teacherOnly.GET("/invitations/active", handlerInvitation.NewGetActiveHandler(uc.Invitation.GetActive))
	teacherOnly.DELETE("/invitations/:id", handlerInvitation.NewRevokeHandler(uc.Invitation.Revoke))
	api.POST("/invitations/redeem", handlerInvitation.NewRedeemHandler(uc.Invitation.Redeem))

	api.POST("/courses/:id/enroll", enrollment.NewEnrollHandler(uc.Enrollment.Enroll))
	api.GET("/courses/:id/students", enrollment.NewListStudentsHandler(uc.Enrollment.ListStudents))

	teacherOnly.POST("/courses/:id/materials", material.NewCreateHandler(uc.Material.Create))
	api.GET("/courses/:id/materials", material.NewListHandler(uc.Material.List))
	api.GET("/materials/:id", material.NewGetHandler(uc.Material.Get))
	teacherOnly.PUT("/materials/:id", material.NewUpdateHandler(uc.Material.Update))
	teacherOnly.DELETE("/materials/:id", material.NewDeleteHandler(uc.Material.Delete))

	teacherOnly.POST("/courses/:id/topics", handlerTopic.NewCreateHandler(uc.Topic.Create))
	api.GET("/courses/:id/topics", handlerTopic.NewListHandler(uc.Topic.List))
	teacherOnly.PUT("/topics/:id", handlerTopic.NewUpdateHandler(uc.Topic.Update))
	teacherOnly.DELETE("/topics/:id", handlerTopic.NewDeleteHandler(uc.Topic.Delete))

	teacherOnly.POST("/topics/:id/exercise-drafts/ai", ai.NewExerciseDraftsHandler(uc.AI.Proxy, uc.Exercise.List))
	teacherOnly.POST("/topics/:id/exercises", exercise.NewCreateHandler(uc.Exercise.Create))
	api.GET("/topics/:id/exercises", exercise.NewListHandler(uc.Exercise.List))
	api.GET("/exercises/:id/statement-image", exercise.NewStatementImageHandler(uc.Exercise.StatementImage))
	teacherOnly.PUT("/exercises/:id", exercise.NewUpdateHandler(uc.Exercise.Update))
	teacherOnly.DELETE("/exercises/:id", exercise.NewDeleteHandler(uc.Exercise.Delete))

	teacherOnly.POST("/courses/:id/practice-sheets", practicesheet.NewCreateHandler(uc.PracticeSheet.Create))
	api.GET("/courses/:id/practice-sheets", practicesheet.NewListHandler(uc.PracticeSheet.List))
	api.GET("/practice-sheets/:id", practicesheet.NewGetHandler(uc.PracticeSheet.Get))
	api.GET("/practice-sheets/:id/exercises/:exerciseId/assistant-media", practicesheet.NewGetAssistantMediaHandler(uc.PracticeSheet.GetAssistantMedia))
	teacherOnly.PUT("/practice-sheets/:id", practicesheet.NewUpdateHandler(uc.PracticeSheet.Update))
	teacherOnly.DELETE("/practice-sheets/:id", practicesheet.NewDeleteHandler(uc.PracticeSheet.Delete))
	api.POST("/practice-sheets/:id/submit", practicesheet.NewSubmitHandler(uc.PracticeSheet.Submit))
	api.POST("/practice-sheets/:id/submit-async", practicesheet.NewSubmitAsyncHandler(uc.PracticeSheet.Submit, submitJobRepo))
	api.GET("/practice-sheets/submit-jobs/:jobId", practicesheet.NewGetSubmitJobHandler(submitJobRepo))

	api.GET("/students/me/progress", studentprogress.NewGetMyProgressHandler(uc.Progress.GetMy))
	api.GET("/students/me/dashboard", studentprogress.NewDashboardHandler(uc.Progress.Dashboard))
	api.GET("/students/me/courses/:id/progress", studentprogress.NewGetCourseProgressHandler(uc.Progress.GetCourse))
	api.GET("/students/me/courses/:id/leaderboard", studentprogress.NewGetCourseLeaderboardHandler(uc.Progress.GetCourseLeaderboard))

	api.GET("/teachers/me/students/:studentId/progress", studentprogress.NewGetStudentProgressHandler(uc.Progress.GetStudentProgress))
	api.GET("/teachers/me/students/:studentId/courses/:courseId/progress", studentprogress.NewGetStudentCourseProgressHandler(uc.Progress.GetStudentCourseProgress))
	api.GET("/teachers/me/students/:studentId/attempts", studentprogress.NewGetStudentAttemptsHandler(uc.Progress.GetStudentAttempts))
	teacherOnly.GET("/teachers/me/students/:studentId/report.pdf", studentreport.NewGeneratePDFHandler(uc.Report.GeneratePDF))

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

	api.GET("/courses/:id/levels", courselevel.NewGetHandler(uc.CourseLevel.Get))

	teacherOnly.POST("/courses/:id/notebooks", handlerNB.NewCreateHandler(uc.Notebook.Create))
	api.GET("/courses/:id/notebooks", handlerNB.NewListHandler(uc.Notebook.List))
	api.GET("/notebooks/:id", handlerNB.NewGetHandler(uc.Notebook.Get))
	teacherOnly.POST("/notebooks/:id/page-drafts/ai", ai.NewNotebookPageDraftsHandler(uc.AI.Proxy, uc.Notebook.Get))
	teacherOnly.PUT("/notebooks/:id", handlerNB.NewUpdateHandler(uc.Notebook.Update))
	teacherOnly.DELETE("/notebooks/:id", handlerNB.NewDeleteHandler(uc.Notebook.Delete))
	teacherOnly.POST("/notebooks/:id/pages", handlerNB.NewAddPageHandler(uc.Notebook.AddPage))
	teacherOnly.PUT("/notebook-pages/:id", handlerNB.NewUpdatePageHandler(uc.Notebook.UpdatePage))
	teacherOnly.DELETE("/notebook-pages/:id", handlerNB.NewDeletePageHandler(uc.Notebook.DeletePage))
	api.POST("/notebook-pages/:id/submit", handlerNB.NewSaveSubmissionHandler(uc.Notebook.SaveSubmission))
	api.POST("/notebook-pages/:id/submit-async", handlerNB.NewSaveSubmissionAsyncHandler(uc.Notebook.SaveSubmission, submitJobRepo))
	api.GET("/notebook-pages/submit-jobs/:jobId", handlerNB.NewGetSubmitJobHandler(submitJobRepo))
	teacherOnly.GET("/notebook-submissions", handlerNB.NewListSubmissionsHandler(uc.Notebook.ListSubmissions))
	teacherOnly.POST("/notebook-submissions/:id/review", handlerNB.NewReviewSubmissionHandler(uc.Notebook.ReviewSubmission))
	teacherOnly.PUT("/notebook-submissions/:id/teacher-review", handlerNB.NewTeacherReviewSubmissionHandler(uc.Notebook.TeacherReview))

	teacherOnly.GET("/subscription-plans", subscription.NewListPlansHandler(uc.Subscription.ListPlans))
	teacherOnly.GET("/teachers/me/subscription", subscription.NewGetMineHandler(uc.Subscription.GetMine))

	teacherOnly.GET("/teachers/me/subscription/checkout-config", subscription.NewCheckoutConfigHandler(config.GetConfigService().ServerConfig.MercadoPagoPublicKey))
	teacherOnly.POST("/teachers/me/subscription", subscription.NewSubscribeHandler(uc.Subscription.Subscribe))

	teacherOnly.POST("/teachers/me/subscription/hosted-checkout", subscription.NewHostedCheckoutHandler(uc.Subscription.HostedCheckout))

	teacherOnly.POST("/teachers/me/subscription/change-plan", subscription.NewChangePlanHandler(uc.Subscription.ChangePlan))
	teacherOnly.POST("/teachers/me/subscription/pause", subscription.NewManageMineHandler(uc.Subscription.ManageMine, ucSubscription.ActionPause))
	teacherOnly.POST("/teachers/me/subscription/resume", subscription.NewManageMineHandler(uc.Subscription.ManageMine, ucSubscription.ActionResume))
	teacherOnly.POST("/teachers/me/subscription/cancel", subscription.NewManageMineHandler(uc.Subscription.ManageMine, ucSubscription.ActionCancel))

	adminOnly.GET("/schools", handlerSchool.NewListHandler(uc.School.List))
	adminOnly.POST("/schools", handlerSchool.NewCreateHandler(uc.School.Create))
	adminOnly.POST("/schools/:id/suspend", handlerSchool.NewSuspendHandler(uc.School.Suspend))
	adminOnly.POST("/schools/:id/close", handlerSchool.NewCloseHandler(uc.School.Close))
	adminOnly.POST("/schools/:id/reopen", handlerSchool.NewReopenHandler(uc.School.Reopen))
	adminOnly.GET("/schools/:id/archive", handlerSchool.NewArchiveHandler(uc.School.Archive))

	api.GET("/schools/mine", handlerSchool.NewMineHandler(uc.School.Mine))

	teacherOnly.PUT("/schools/:id", handlerSchool.NewUpdateHandler(uc.School.Update))
	teacherOnly.GET("/schools/:id/members", handlerSchool.NewListMembersHandler(uc.School.ListMembers))
	teacherOnly.POST("/schools/:id/members", handlerSchool.NewAddMemberHandler(uc.School.AddMember))
	teacherOnly.DELETE("/schools/:id/members/:userId", handlerSchool.NewRemoveMemberHandler(uc.School.RemoveMember))

	teacherOnly.GET("/teachers/me/subscription/downgrade", subscription.NewDowngradePreviewHandler(uc.Subscription.DowngradePreview))
	teacherOnly.POST("/teachers/me/subscription/downgrade", subscription.NewDowngradeApplyHandler(uc.Subscription.DowngradeApply))
	teacherOnly.POST("/teachers/me/students/:studentId/reactivate", subscription.NewReactivateStudentHandler(uc.Subscription.ReactivateStudent))

	adminOnly.POST("/subscription-plans", subscription.NewCreatePlanHandler(uc.Subscription.CreatePlan))
	adminOnly.PUT("/subscription-plans/:id", subscription.NewUpdatePlanHandler(uc.Subscription.UpdatePlan))
	adminOnly.DELETE("/subscription-plans/:id", subscription.NewDeactivatePlanHandler(uc.Subscription.DeactivatePlan))

	teacherOnly.GET("/attempt-reviews", handlerReview.NewListHandler(uc.AttemptReview.List))
	teacherOnly.POST("/attempt-reviews/:id", handlerReview.NewReviewHandler(uc.AttemptReview.Review))
	teacherOnly.GET("/attempt-reviews/:id/statement-image", handlerReview.NewStatementImageHandler(uc.AttemptReview.StatementImage))

	api.POST("/uploads", handlerUpload.NewHandler(uc.Upload.Upload))

	api.GET("/notifications", handlerNotification.NewListHandler(uc.Notification.List))
	api.POST("/notifications/:id/read", handlerNotification.NewMarkReadHandler(uc.Notification.MarkRead))
	api.POST("/notifications/read-all", handlerNotification.NewMarkAllReadHandler(uc.Notification.MarkAllRead))
	api.DELETE("/notifications/:id", handlerNotification.NewDeleteHandler(uc.Notification.Delete))

	api.GET("/learning-strategies", handlerLS.NewListHandler(uc.LearningStrategy.List))
	api.GET("/learning-strategies/:id", handlerLS.NewGetHandler(uc.LearningStrategy.Get))
	adminOnly.POST("/learning-strategies", handlerLS.NewCreateHandler(uc.LearningStrategy.Create))
	adminOnly.PUT("/learning-strategies/:id", handlerLS.NewUpdateHandler(uc.LearningStrategy.Update))
	adminOnly.DELETE("/learning-strategies/:id", handlerLS.NewDeleteHandler(uc.LearningStrategy.Delete))

	api.GET("/courses/:id/strategies", handlerLS.NewListByCourseHandler(uc.LearningStrategy.ListByCourse))
	teacherOnly.POST("/courses/:id/strategies", handlerLS.NewAssignToCourseHandler(uc.LearningStrategy.AssignToCourse))
	teacherOnly.DELETE("/course-learning-strategies/:id", handlerLS.NewUnassignFromCourseHandler(uc.LearningStrategy.UnassignFromCourse))

	teacherOnly.GET("/students/:studentId/courses/:courseId/progress", handlerCP.NewGetForStudentHandler(uc.CourseProgress.GetForStudent))
	teacherOnly.GET("/students/:studentId/progress", handlerCP.NewListForStudentHandler(uc.CourseProgress.ListForStudent))
}
