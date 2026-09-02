package repositories

import (
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources"
	aiconversation "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/ai_conversation"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	coursecuriosities "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course_curiosities"
	courseprogress "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course_progress"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/exercise"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/grade"
	learningstrategy "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/learning_strategy"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/material"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notification"
	practicesheet "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/practice_sheet"
	studentattempt "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_attempt"
	studentinvitation "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_invitation"
	studentprogress "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_progress"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/subject"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	teacherstudentassignment "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/teacher_student_assignment"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/topic"
	userprofile "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/user_profile"
)

type Repositories struct {
	UserProfile              userprofile.Repository
	Grade                    grade.Repository
	Subject                  subject.Repository
	TeacherStudentAssignment teacherstudentassignment.Repository
	Course                   course.Repository
	CourseCuriosities        coursecuriosities.Repository
	Topic                    topic.Repository
	Exercise                 exercise.Repository
	Material                 material.Repository
	PracticeSheet            practicesheet.Repository
	Enrollment               enrollment.Repository
	StudentAttempt           studentattempt.Repository
	StudentProgress          studentprogress.Repository
	StudentInvitation        studentinvitation.Repository
	AIConversation           aiconversation.Repository
	LearningStrategy         learningstrategy.Repository
	Notebook                 notebook.Repository
	CourseProgress           courseprogress.Repository
	SubmitJob                submitjob.Repository
	Notification             notification.Repository
}

type Factory func() *Repositories

func NewFactory(ds *datasources.Datasources) func() *Repositories {
	return func() *Repositories {
		return &Repositories{
			UserProfile:              userprofile.NewRepository(ds.DB),
			Grade:                    grade.NewRepository(ds.DB),
			Subject:                  subject.NewRepository(ds.DB),
			TeacherStudentAssignment: teacherstudentassignment.NewRepository(ds.DB),
			Course:                   course.NewRepository(ds.DB),
			CourseCuriosities:        coursecuriosities.NewRepository(ds.DB),
			Topic:                    topic.NewRepository(ds.DB),
			Exercise:                 exercise.NewRepository(ds.DB),
			Material:                 material.NewRepository(ds.DB),
			PracticeSheet:            practicesheet.NewRepository(ds.DB),
			Enrollment:               enrollment.NewRepository(ds.DB),
			StudentAttempt:           studentattempt.NewRepository(ds.DB),
			StudentProgress:          studentprogress.NewRepository(ds.DB),
			StudentInvitation:        studentinvitation.NewRepository(ds.DB),
			AIConversation:           aiconversation.NewRepository(ds.DB),
			LearningStrategy:         learningstrategy.NewRepository(ds.DB),
			Notebook:                 notebook.NewRepository(ds.DB),
			CourseProgress:           courseprogress.NewRepository(ds.DB),
			SubmitJob:                submitjob.NewRepository(ds.DB),
			Notification:             notification.NewRepository(ds.DB),
		}
	}
}
