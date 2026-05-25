package usecases

import (
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"
	ucLevel "github.com/tapiaw38/practiq-be/internal/usecases/course_level"
	ucEnrollment "github.com/tapiaw38/practiq-be/internal/usecases/enrollment"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"
	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
	ucPracticeSheet "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
	ucProgress "github.com/tapiaw38/practiq-be/internal/usecases/student_progress"
	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"
	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"
	ucTopic "github.com/tapiaw38/practiq-be/internal/usecases/topic"
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
	GetMy     ucProgress.GetMyProgressUsecase
	GetCourse ucProgress.GetCourseProgressUsecase
}

type AIUsecases struct {
	CreateConversation ucAI.CreateConversationUsecase
	GetMessages        ucAI.GetMessagesUsecase
	Help               ucAI.HelpUsecase
}

type ProfileUsecases struct {
	Sync                  ucProfile.SyncUsecase
	Get                   ucProfile.GetUsecase
	UpdateAssistantConfig ucProfile.UpdateAssistantConfigUsecase
	UpdateAcademicStatus  ucProfile.UpdateAcademicStatusUsecase
}

type NotebookUsecases struct {
	Create         ucNB.CreateUsecase
	List           ucNB.ListUsecase
	Get            ucNB.GetUsecase
	Update         ucNB.UpdateUsecase
	Delete         ucNB.DeleteUsecase
	AddPage        ucNB.AddPageUsecase
	UpdatePage     ucNB.UpdatePageUsecase
	SaveSubmission ucNB.SaveSubmissionUsecase
}

type CourseLevelUsecases struct {
	Get ucLevel.GetUsecase
}

type Usecases struct {
	Course        CourseUsecases
	Grade         GradeUsecases
	Subject       SubjectUsecases
	Assignment    AssignmentUsecases
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
		Grade: GradeUsecases{
			Create:         ucGrade.NewCreateUsecase(factory),
			List:           ucGrade.NewListUsecase(factory),
			Update:         ucGrade.NewUpdateUsecase(factory),
			Delete:         ucGrade.NewDeleteUsecase(factory),
			AssignMember:   ucGrade.NewAssignMemberUsecase(factory),
			ListMembers:    ucGrade.NewListMembersUsecase(factory),
			RemoveMember:   ucGrade.NewRemoveMemberUsecase(factory),
			ListUserGrades: ucGrade.NewListUserGradesUsecase(factory),
		},
		Subject: SubjectUsecases{
			Create: ucSubject.NewCreateUsecase(factory),
			List:   ucSubject.NewListUsecase(factory),
			Update: ucSubject.NewUpdateUsecase(factory),
			Delete: ucSubject.NewDeleteUsecase(factory),
		},
		Assignment: AssignmentUsecases{
			Assign:       ucAssignment.NewAssignUsecase(factory),
			Unassign:     ucAssignment.NewUnassignUsecase(factory),
			ListStudents: ucAssignment.NewListStudentsUsecase(factory),
			ListTeachers: ucAssignment.NewListTeachersUsecase(factory),
		},
		Topic: TopicUsecases{
			Create: ucTopic.NewCreateUsecase(factory),
			List:   ucTopic.NewListUsecase(factory),
			Update: ucTopic.NewUpdateUsecase(factory),
			Delete: ucTopic.NewDeleteUsecase(factory),
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
			Update: ucMaterial.NewUpdateUsecase(factory),
			Delete: ucMaterial.NewDeleteUsecase(factory),
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
			Update: ucPracticeSheet.NewUpdateUsecase(factory),
			Delete: ucPracticeSheet.NewDeleteUsecase(factory),
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
			Sync:                  ucProfile.NewSyncUsecase(factory),
			Get:                   ucProfile.NewGetUsecase(factory),
			UpdateAssistantConfig: ucProfile.NewUpdateAssistantConfigUsecase(factory),
			UpdateAcademicStatus:  ucProfile.NewUpdateAcademicStatusUsecase(factory),
		},
		Notebook: NotebookUsecases{
			Create:         ucNB.NewCreateUsecase(factory),
			List:           ucNB.NewListUsecase(factory),
			Get:            ucNB.NewGetUsecase(factory),
			Update:         ucNB.NewUpdateUsecase(factory),
			Delete:         ucNB.NewDeleteUsecase(factory),
			AddPage:        ucNB.NewAddPageUsecase(factory),
			UpdatePage:     ucNB.NewUpdatePageUsecase(factory),
			SaveSubmission: ucNB.NewSaveSubmissionUsecase(factory),
		},
		CourseLevel: CourseLevelUsecases{
			Get: ucLevel.NewGetUsecase(factory),
		},
	}
}
