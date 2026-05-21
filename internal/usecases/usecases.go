package usecases

import (
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"
	ucLevel "github.com/tapiaw38/practiq-be/internal/usecases/course_level"
	ucEnrollment "github.com/tapiaw38/practiq-be/internal/usecases/enrollment"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"
	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
	ucPracticeSheet "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"
	ucTopic "github.com/tapiaw38/practiq-be/internal/usecases/topic"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type CourseUsecases struct {
	Create ucCourse.CreateUsecase
	List   ucCourse.ListUsecase
	Get    ucCourse.GetUsecase
	Update ucCourse.UpdateUsecase
	Delete ucCourse.DeleteUsecase
}

type TopicUsecases struct {
	Create ucTopic.CreateUsecase
	List   ucTopic.ListUsecase
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
}

type ProgressUsecases struct {
	GetMy     ucProgress.GetMyProgressUsecase
	GetCourse ucProgress.GetCourseProgressUsecase
}

type AIUsecases struct {
	CreateConversation ucAI.CreateConversationUsecase
	GetMessages        ucAI.GetMessagesUsecase
	Help               ucAI.HelpUsecase
}

type ProfileUsecases struct {
	Sync ucProfile.SyncUsecase
	Get  ucProfile.GetUsecase
}

type NotebookUsecases struct {
	Create         ucNB.CreateUsecase
	List           ucNB.ListUsecase
	Get            ucNB.GetUsecase
	AddPage        ucNB.AddPageUsecase
	UpdatePage     ucNB.UpdatePageUsecase
	SaveSubmission ucNB.SaveSubmissionUsecase
}

type CourseLevelUsecases struct {
	Get ucLevel.GetUsecase
}

type Usecases struct {
	Course        CourseUsecases
	Topic         TopicUsecases
	Exercise      ExerciseUsecases
	Material      MaterialUsecases
	Enrollment    EnrollmentUsecases
	PracticeSheet PracticeSheetUsecases
	Progress      ProgressUsecases
	AI            AIUsecases
	Profile       ProfileUsecases
	Notebook      NotebookUsecases
	CourseLevel   CourseLevelUsecases
}

func NewUsecases(factory appcontext.Factory) *Usecases {
	return &Usecases{
		Course: CourseUsecases{
			Create: ucCourse.NewCreateUsecase(factory),
			List:   ucCourse.NewListUsecase(factory),
			Get:    ucCourse.NewGetUsecase(factory),
			Update: ucCourse.NewUpdateUsecase(factory),
			Delete: ucCourse.NewDeleteUsecase(factory),
		},
		Topic: TopicUsecases{
			Create: ucTopic.NewCreateUsecase(factory),
			List:   ucTopic.NewListUsecase(factory),
		},
		Exercise: ExerciseUsecases{
			Create: ucExercise.NewCreateUsecase(factory),
			List:   ucExercise.NewListUsecase(factory),
			Update: ucExercise.NewUpdateUsecase(factory),
			Delete: ucExercise.NewDeleteUsecase(factory),
		},
		Material: MaterialUsecases{
			Create: ucMaterial.NewCreateUsecase(factory),
			List:   ucMaterial.NewListUsecase(factory),
		},
		Enrollment: EnrollmentUsecases{
			Enroll:       ucEnrollment.NewEnrollUsecase(factory),
			ListStudents: ucEnrollment.NewListStudentsUsecase(factory),
		},
		PracticeSheet: PracticeSheetUsecases{
			Create: ucPracticeSheet.NewCreateUsecase(factory),
			List:   ucPracticeSheet.NewListUsecase(factory),
			Get:    ucPracticeSheet.NewGetUsecase(factory),
			Submit: ucPracticeSheet.NewSubmitUsecase(factory),
		},
		Progress: ProgressUsecases{
			GetMy:     ucProgress.NewGetMyProgressUsecase(factory),
			GetCourse: ucProgress.NewGetCourseProgressUsecase(factory),
		},
		AI: AIUsecases{
			CreateConversation: ucAI.NewCreateConversationUsecase(factory),
			GetMessages:        ucAI.NewGetMessagesUsecase(factory),
			Help:               ucAI.NewHelpUsecase(factory),
		},
		Profile: ProfileUsecases{
			Sync: ucProfile.NewSyncUsecase(factory),
			Get:  ucProfile.NewGetUsecase(factory),
		},
		Notebook: NotebookUsecases{
			Create:         ucNB.NewCreateUsecase(factory),
			List:           ucNB.NewListUsecase(factory),
			Get:            ucNB.NewGetUsecase(factory),
			AddPage:        ucNB.NewAddPageUsecase(factory),
			UpdatePage:     ucNB.NewUpdatePageUsecase(factory),
			SaveSubmission: ucNB.NewSaveSubmissionUsecase(factory),
		},
		CourseLevel: CourseLevelUsecases{
			Get: ucLevel.NewGetUsecase(factory),
		},
	}
}
