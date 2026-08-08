package usecases

import (
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
	ucAttemptReview "github.com/tapiaw38/practiq-be/internal/usecases/attempt_review"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"
	ucLevel "github.com/tapiaw38/practiq-be/internal/usecases/course_level"
	ucCourseProgress "github.com/tapiaw38/practiq-be/internal/usecases/course_progress"
	ucEnrollment "github.com/tapiaw38/practiq-be/internal/usecases/enrollment"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"
	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"
	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
	ucNotification "github.com/tapiaw38/practiq-be/internal/usecases/notification"
	ucPracticeSheet "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"
	ucReport "github.com/tapiaw38/practiq-be/internal/usecases/student_report"
	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"
	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"
	ucTopic "github.com/tapiaw38/practiq-be/internal/usecases/topic"
	ucUpload "github.com/tapiaw38/practiq-be/internal/usecases/upload"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"
)

type CourseUsecases struct {
	Create ucCourse.CreateUsecase
	List   ucCourse.ListUsecase
	Get    ucCourse.GetUsecase
	Update ucCourse.UpdateUsecase
	Delete ucCourse.DeleteUsecase
}

type GradeUsecases struct {
	Create         ucGrade.CreateUsecase
	List           ucGrade.ListUsecase
	Update         ucGrade.UpdateUsecase
	Delete         ucGrade.DeleteUsecase
	AssignMember   ucGrade.AssignMemberUsecase
	ListMembers    ucGrade.ListMembersUsecase
	RemoveMember   ucGrade.RemoveMemberUsecase
	ListUserGrades ucGrade.ListUserGradesUsecase
}

type SubjectUsecases struct {
	Create ucSubject.CreateUsecase
	List   ucSubject.ListUsecase
	Update ucSubject.UpdateUsecase
	Delete ucSubject.DeleteUsecase
}

type AssignmentUsecases struct {
	Assign       ucAssignment.AssignUsecase
	Unassign     ucAssignment.UnassignUsecase
	ListStudents ucAssignment.ListStudentsUsecase
	ListTeachers ucAssignment.ListTeachersUsecase
}

type TopicUsecases struct {
	Create ucTopic.CreateUsecase
	List   ucTopic.ListUsecase
	Update ucTopic.UpdateUsecase
	Delete ucTopic.DeleteUsecase
}

type ExerciseUsecases struct {
	Create ucExercise.CreateUsecase
	List   ucExercise.ListUsecase
	Update ucExercise.UpdateUsecase
	Delete ucExercise.DeleteUsecase
}

type MaterialUsecases struct {
	Create ucMaterial.CreateUsecase
	List   ucMaterial.ListUsecase
	Update ucMaterial.UpdateUsecase
	Delete ucMaterial.DeleteUsecase
}

type EnrollmentUsecases struct {
	Enroll       ucEnrollment.EnrollUsecase
	ListStudents ucEnrollment.ListStudentsUsecase
}

type PracticeSheetUsecases struct {
	Create ucPracticeSheet.CreateUsecase
	List   ucPracticeSheet.ListUsecase
	Get    ucPracticeSheet.GetUsecase
	Submit ucPracticeSheet.SubmitUsecase
	Update ucPracticeSheet.UpdateUsecase
	Delete ucPracticeSheet.DeleteUsecase
}

type ProgressUsecases struct {
	GetMy                    ucProgress.GetMyProgressUsecase
	GetCourse                ucProgress.GetCourseProgressUsecase
	GetStudentProgress       ucProgress.GetStudentProgressUsecase
	GetStudentCourseProgress ucProgress.GetStudentCourseProgressUsecase
	GetStudentAttempts       ucProgress.GetStudentAttemptsUsecase
}

type AIUsecases struct {
	CreateConversation  ucAI.CreateConversationUsecase
	GetMessages         ucAI.GetMessagesUsecase
	Help                ucAI.HelpUsecase
	Proxy               ucAI.ProxyUsecase
	GenerateCuriosities ucAI.GenerateCuriositiesUsecase
}

type ProfileUsecases struct {
	Sync                  ucProfile.SyncUsecase
	Get                   ucProfile.GetUsecase
	UpdateAssistantConfig ucProfile.UpdateAssistantConfigUsecase
	UpdateAcademicStatus  ucProfile.UpdateAcademicStatusUsecase
}

type NotebookUsecases struct {
	Create           ucNB.CreateUsecase
	List             ucNB.ListUsecase
	Get              ucNB.GetUsecase
	Update           ucNB.UpdateUsecase
	Delete           ucNB.DeleteUsecase
	AddPage          ucNB.AddPageUsecase
	UpdatePage       ucNB.UpdatePageUsecase
	SaveSubmission   ucNB.SaveSubmissionUsecase
	ListSubmissions  ucNB.ListSubmissionsUsecase
	ReviewSubmission ucNB.ReviewSubmissionUsecase
	TeacherReview    ucNB.TeacherReviewSubmissionUsecase
}

type CourseLevelUsecases struct {
	Get ucLevel.GetUsecase
}

type LearningStrategyUsecases struct {
	List               ucLS.ListUsecase
	Get                ucLS.GetUsecase
	Create             ucLS.CreateUsecase
	Update             ucLS.UpdateUsecase
	Delete             ucLS.DeleteUsecase
	ListByCourse       ucLS.ListByCourseUsecase
	AssignToCourse     ucLS.AssignToCourseUsecase
	UnassignFromCourse ucLS.UnassignFromCourseUsecase
}

type CourseProgressUsecases struct {
	GetForStudent  ucCourseProgress.GetForStudentUsecase
	ListForStudent ucCourseProgress.ListForStudentUsecase
}

type AttemptReviewUsecases struct {
	List   ucAttemptReview.ListUsecase
	Review ucAttemptReview.ReviewUsecase
}

type UploadUsecases struct {
	Upload ucUpload.Usecase
}

type NotificationUsecases struct {
	List        ucNotification.ListUsecase
	MarkRead    ucNotification.MarkReadUsecase
	MarkAllRead ucNotification.MarkAllReadUsecase
}

type ReportUsecases struct {
	GeneratePDF ucReport.GeneratePDFUsecase
}

type Usecases struct {
	Course           CourseUsecases
	Grade            GradeUsecases
	Subject          SubjectUsecases
	Assignment       AssignmentUsecases
	Topic            TopicUsecases
	Exercise         ExerciseUsecases
	Material         MaterialUsecases
	Enrollment       EnrollmentUsecases
	PracticeSheet    PracticeSheetUsecases
	Progress         ProgressUsecases
	AI               AIUsecases
	Profile          ProfileUsecases
	Notebook         NotebookUsecases
	CourseLevel      CourseLevelUsecases
	LearningStrategy LearningStrategyUsecases
	CourseProgress   CourseProgressUsecases
	Report           ReportUsecases
	Notification     NotificationUsecases
	Upload           UploadUsecases
	AttemptReview    AttemptReviewUsecases
}

func NewUsecases(contextFactory appcontext.Factory) *Usecases {
	return &Usecases{
		Course: CourseUsecases{
			Create: ucCourse.NewCreateUsecase(contextFactory),
			List:   ucCourse.NewListUsecase(contextFactory),
			Get:    ucCourse.NewGetUsecase(contextFactory),
			Update: ucCourse.NewUpdateUsecase(contextFactory),
			Delete: ucCourse.NewDeleteUsecase(contextFactory),
		},
		AttemptReview: AttemptReviewUsecases{
			List:   ucAttemptReview.NewListUsecase(contextFactory),
			Review: ucAttemptReview.NewReviewUsecase(contextFactory),
		},
		Upload: UploadUsecases{
			Upload: ucUpload.NewUsecase(contextFactory),
		},
		Notification: NotificationUsecases{
			List:        ucNotification.NewListUsecase(contextFactory),
			MarkRead:    ucNotification.NewMarkReadUsecase(contextFactory),
			MarkAllRead: ucNotification.NewMarkAllReadUsecase(contextFactory),
		},
		Grade: GradeUsecases{
			Create:         ucGrade.NewCreateUsecase(contextFactory),
			List:           ucGrade.NewListUsecase(contextFactory),
			Update:         ucGrade.NewUpdateUsecase(contextFactory),
			Delete:         ucGrade.NewDeleteUsecase(contextFactory),
			AssignMember:   ucGrade.NewAssignMemberUsecase(contextFactory),
			ListMembers:    ucGrade.NewListMembersUsecase(contextFactory),
			RemoveMember:   ucGrade.NewRemoveMemberUsecase(contextFactory),
			ListUserGrades: ucGrade.NewListUserGradesUsecase(contextFactory),
		},
		Subject: SubjectUsecases{
			Create: ucSubject.NewCreateUsecase(contextFactory),
			List:   ucSubject.NewListUsecase(contextFactory),
			Update: ucSubject.NewUpdateUsecase(contextFactory),
			Delete: ucSubject.NewDeleteUsecase(contextFactory),
		},
		Assignment: AssignmentUsecases{
			Assign:       ucAssignment.NewAssignUsecase(contextFactory),
			Unassign:     ucAssignment.NewUnassignUsecase(contextFactory),
			ListStudents: ucAssignment.NewListStudentsUsecase(contextFactory),
			ListTeachers: ucAssignment.NewListTeachersUsecase(contextFactory),
		},
		Topic: TopicUsecases{
			Create: ucTopic.NewCreateUsecase(contextFactory),
			List:   ucTopic.NewListUsecase(contextFactory),
			Update: ucTopic.NewUpdateUsecase(contextFactory),
			Delete: ucTopic.NewDeleteUsecase(contextFactory),
		},
		Exercise: ExerciseUsecases{
			Create: ucExercise.NewCreateUsecase(contextFactory),
			List:   ucExercise.NewListUsecase(contextFactory),
			Update: ucExercise.NewUpdateUsecase(contextFactory),
			Delete: ucExercise.NewDeleteUsecase(contextFactory),
		},
		Material: MaterialUsecases{
			Create: ucMaterial.NewCreateUsecase(contextFactory),
			List:   ucMaterial.NewListUsecase(contextFactory),
			Update: ucMaterial.NewUpdateUsecase(contextFactory),
			Delete: ucMaterial.NewDeleteUsecase(contextFactory),
		},
		Enrollment: EnrollmentUsecases{
			Enroll:       ucEnrollment.NewEnrollUsecase(contextFactory),
			ListStudents: ucEnrollment.NewListStudentsUsecase(contextFactory),
		},
		PracticeSheet: PracticeSheetUsecases{
			Create: ucPracticeSheet.NewCreateUsecase(contextFactory),
			List:   ucPracticeSheet.NewListUsecase(contextFactory),
			Get:    ucPracticeSheet.NewGetUsecase(contextFactory),
			Submit: ucPracticeSheet.NewSubmitUsecase(contextFactory),
			Update: ucPracticeSheet.NewUpdateUsecase(contextFactory),
			Delete: ucPracticeSheet.NewDeleteUsecase(contextFactory),
		},
		Progress: ProgressUsecases{
			GetMy:                    ucProgress.NewGetMyProgressUsecase(contextFactory),
			GetCourse:                ucProgress.NewGetCourseProgressUsecase(contextFactory),
			GetStudentProgress:       ucProgress.NewGetStudentProgressUsecase(contextFactory),
			GetStudentCourseProgress: ucProgress.NewGetStudentCourseProgressUsecase(contextFactory),
			GetStudentAttempts:       ucProgress.NewGetStudentAttemptsUsecase(contextFactory),
		},
		AI: AIUsecases{
			CreateConversation:  ucAI.NewCreateConversationUsecase(contextFactory),
			GetMessages:         ucAI.NewGetMessagesUsecase(contextFactory),
			Help:                ucAI.NewHelpUsecase(contextFactory),
			Proxy:               ucAI.NewProxyUsecase(contextFactory),
			GenerateCuriosities: ucAI.NewGenerateCuriositiesUsecase(contextFactory),
		},
		Profile: ProfileUsecases{
			Sync:                  ucProfile.NewSyncUsecase(contextFactory),
			Get:                   ucProfile.NewGetUsecase(contextFactory),
			UpdateAssistantConfig: ucProfile.NewUpdateAssistantConfigUsecase(contextFactory),
			UpdateAcademicStatus:  ucProfile.NewUpdateAcademicStatusUsecase(contextFactory),
		},
		Notebook: NotebookUsecases{
			Create:           ucNB.NewCreateUsecase(contextFactory),
			List:             ucNB.NewListUsecase(contextFactory),
			Get:              ucNB.NewGetUsecase(contextFactory),
			Update:           ucNB.NewUpdateUsecase(contextFactory),
			Delete:           ucNB.NewDeleteUsecase(contextFactory),
			AddPage:          ucNB.NewAddPageUsecase(contextFactory),
			UpdatePage:       ucNB.NewUpdatePageUsecase(contextFactory),
			SaveSubmission:   ucNB.NewSaveSubmissionUsecase(contextFactory),
			ListSubmissions:  ucNB.NewListSubmissionsUsecase(contextFactory),
			ReviewSubmission: ucNB.NewReviewSubmissionUsecase(contextFactory),
			TeacherReview:    ucNB.NewTeacherReviewSubmissionUsecase(contextFactory),
		},
		CourseLevel: CourseLevelUsecases{
			Get: ucLevel.NewGetUsecase(contextFactory),
		},
		LearningStrategy: LearningStrategyUsecases{
			List:               ucLS.NewListUsecase(contextFactory),
			Get:                ucLS.NewGetUsecase(contextFactory),
			Create:             ucLS.NewCreateUsecase(contextFactory),
			Update:             ucLS.NewUpdateUsecase(contextFactory),
			Delete:             ucLS.NewDeleteUsecase(contextFactory),
			ListByCourse:       ucLS.NewListByCourseUsecase(contextFactory),
			AssignToCourse:     ucLS.NewAssignToCourseUsecase(contextFactory),
			UnassignFromCourse: ucLS.NewUnassignFromCourseUsecase(contextFactory),
		},
		CourseProgress: CourseProgressUsecases{
			GetForStudent:  ucCourseProgress.NewGetForStudentUsecase(contextFactory),
			ListForStudent: ucCourseProgress.NewListForStudentUsecase(contextFactory),
		},
		Report: ReportUsecases{
			GeneratePDF: ucReport.NewGeneratePDFUsecase(contextFactory),
		},
	}
}
