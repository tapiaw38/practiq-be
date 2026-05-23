package repositories

import (
	"database/sql"

	aiconversation "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/ai_conversation"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	courseprogress "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course_progress"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/exercise"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/grade"
	learningstrategy "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/learning_strategy"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/material"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
	practicesheet "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/practice_sheet"
	studentattempt "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_attempt"
	studentprogress "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_progress"
	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/subject"
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
	Topic                    topic.Repository
	Exercise                 exercise.Repository
	Material                 material.Repository
	PracticeSheet            practicesheet.Repository
	Enrollment               enrollment.Repository
	StudentAttempt           studentattempt.Repository
	StudentProgress          studentprogress.Repository
	AIConversation           aiconversation.Repository
	LearningStrategy         learningstrategy.Repository
	Notebook                 notebook.Repository
	CourseProgress           courseprogress.Repository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserProfile:              userprofile.NewRepository(db),
		Grade:                    grade.NewRepository(db),
		Subject:                  subject.NewRepository(db),
		TeacherStudentAssignment: teacherstudentassignment.NewRepository(db),
		Course:                   course.NewRepository(db),
		Topic:                    topic.NewRepository(db),
		Exercise:                 exercise.NewRepository(db),
		Material:                 material.NewRepository(db),
		PracticeSheet:            practicesheet.NewRepository(db),
		Enrollment:               enrollment.NewRepository(db),
		StudentAttempt:           studentattempt.NewRepository(db),
		StudentProgress:          studentprogress.NewRepository(db),
		AIConversation:           aiconversation.NewRepository(db),
		LearningStrategy:         learningstrategy.NewRepository(db),
		Notebook:                 notebook.NewRepository(db),
		CourseProgress:           courseprogress.NewRepository(db),
	}
}
