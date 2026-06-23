package mappings

import "net/http"

var (
	// Course errors
	CourseCreateError = ErrorDetails{
		InternalCode: "course:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create course",
	}
	CourseGetError = ErrorDetails{
		InternalCode: "course:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get course",
	}
	CourseNotFoundError = ErrorDetails{
		InternalCode: "course:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "course not found",
	}
	CourseUpdateError = ErrorDetails{
		InternalCode: "course:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update course",
	}
	CourseDeleteError = ErrorDetails{
		InternalCode: "course:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete course",
	}
	CourseListError = ErrorDetails{
		InternalCode: "course:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list courses",
	}

	// Grade errors
	GradeCreateError = ErrorDetails{
		InternalCode: "grade:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create grade",
	}
	GradeListError = ErrorDetails{
		InternalCode: "grade:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list grades",
	}
	GradeGetError = ErrorDetails{
		InternalCode: "grade:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get grade",
	}
	GradeNotFoundError = ErrorDetails{
		InternalCode: "grade:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "grade not found",
	}
	GradeUpdateError = ErrorDetails{
		InternalCode: "grade:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update grade",
	}
	GradeDeleteError = ErrorDetails{
		InternalCode: "grade:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete grade",
	}
	GradeAssignMemberError = ErrorDetails{
		InternalCode: "grade:assign-member:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to assign member to grade",
	}
	GradeListMembersError = ErrorDetails{
		InternalCode: "grade:list-members:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list grade members",
	}

	// Subject errors
	SubjectCreateError = ErrorDetails{
		InternalCode: "subject:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create subject",
	}
	SubjectUpdateError = ErrorDetails{
		InternalCode: "subject:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update subject",
	}
	SubjectDeleteError = ErrorDetails{
		InternalCode: "subject:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete subject",
	}
	SubjectListError = ErrorDetails{
		InternalCode: "subject:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list subjects",
	}
	SubjectGetError = ErrorDetails{
		InternalCode: "subject:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get subject",
	}
	SubjectNotFoundError = ErrorDetails{
		InternalCode: "subject:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "subject not found",
	}

	// Teacher/student assignment errors
	AssignmentCreateError = ErrorDetails{
		InternalCode: "assignment:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to assign student to teacher",
	}
	AssignmentDeleteError = ErrorDetails{
		InternalCode: "assignment:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to unassign student from teacher",
	}
	AssignmentListError = ErrorDetails{
		InternalCode: "assignment:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list assignments",
	}

	// Topic errors
	TopicCreateError = ErrorDetails{
		InternalCode: "topic:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create topic",
	}
	TopicListError = ErrorDetails{
		InternalCode: "topic:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list topics",
	}
	TopicGetError = ErrorDetails{
		InternalCode: "topic:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get topic",
	}
	TopicNotFoundError = ErrorDetails{
		InternalCode: "topic:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "topic not found",
	}
	TopicUpdateError = ErrorDetails{
		InternalCode: "topic:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update topic",
	}
	TopicDeleteError = ErrorDetails{
		InternalCode: "topic:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete topic",
	}

	// Exercise errors
	ExerciseCreateError = ErrorDetails{
		InternalCode: "exercise:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create exercise",
	}
	ExerciseListError = ErrorDetails{
		InternalCode: "exercise:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list exercises",
	}
	ExerciseUpdateError = ErrorDetails{
		InternalCode: "exercise:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update exercise",
	}
	ExerciseDeleteError = ErrorDetails{
		InternalCode: "exercise:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete exercise",
	}
	ExerciseNotFoundError = ErrorDetails{
		InternalCode: "exercise:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "exercise not found",
	}

	// Material errors
	MaterialCreateError = ErrorDetails{
		InternalCode: "material:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create material",
	}
	MaterialListError = ErrorDetails{
		InternalCode: "material:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list materials",
	}
	MaterialGetError = ErrorDetails{
		InternalCode: "material:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get material",
	}
	MaterialNotFoundError = ErrorDetails{
		InternalCode: "material:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "material not found",
	}
	MaterialUpdateError = ErrorDetails{
		InternalCode: "material:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update material",
	}
	MaterialDeleteError = ErrorDetails{
		InternalCode: "material:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete material",
	}

	// Practice sheet errors
	PracticeSheetCreateError = ErrorDetails{
		InternalCode: "practice-sheet:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create practice sheet",
	}
	PracticeSheetListError = ErrorDetails{
		InternalCode: "practice-sheet:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list practice sheets",
	}
	PracticeSheetGetError = ErrorDetails{
		InternalCode: "practice-sheet:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get practice sheet",
	}
	PracticeSheetNotFoundError = ErrorDetails{
		InternalCode: "practice-sheet:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "practice sheet not found",
	}
	PracticeSheetUpdateError = ErrorDetails{
		InternalCode: "practice-sheet:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update practice sheet",
	}
	PracticeSheetDeleteError = ErrorDetails{
		InternalCode: "practice-sheet:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete practice sheet",
	}
	PracticeSheetSubmitError = ErrorDetails{
		InternalCode: "practice-sheet:submit:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to submit practice sheet",
	}

	// Enrollment errors
	EnrollmentCreateError = ErrorDetails{
		InternalCode: "enrollment:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to enroll student",
	}
	EnrollmentListError = ErrorDetails{
		InternalCode: "enrollment:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list students",
	}
	EnrollmentAlreadyExistsError = ErrorDetails{
		InternalCode: "enrollment:create:already-exists",
		StatusCode: http.StatusConflict,
		Message:    "student already enrolled in this course",
	}

	// Progress errors
	ProgressGetError = ErrorDetails{
		InternalCode: "progress:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get progress",
	}
	AttemptGetError = ErrorDetails{
		InternalCode: "attempt:get:error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get attempt",
	}

	// AI errors
	AIConversationCreateError = ErrorDetails{
		InternalCode: "ai:conversation:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create conversation",
	}
	AIConversationGetError = ErrorDetails{
		InternalCode: "ai:conversation:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get conversation",
	}
	AIMessageListError = ErrorDetails{
		InternalCode: "ai:messages:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list messages",
	}
	AIHelpError = ErrorDetails{
		InternalCode: "ai:help:error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get AI help",
	}
	AICuriositiesError = ErrorDetails{
		InternalCode: "ai:curiosities:error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to generate curiosities",
	}

	// Profile errors
	ProfileSyncError = ErrorDetails{
		InternalCode: "profile:sync:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to sync profile",
	}
	ProfileGetError = ErrorDetails{
		InternalCode: "profile:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get profile",
	}
	ProfileUpdateError = ErrorDetails{
		InternalCode: "profile:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update profile",
	}

	// Notebook errors
	NotebookUpdateError = ErrorDetails{
		InternalCode: "notebook:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update notebook",
	}
	NotebookDeleteError = ErrorDetails{
		InternalCode: "notebook:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete notebook",
	}
	NotebookGetError = ErrorDetails{
		InternalCode: "notebook:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get notebook",
	}
	NotebookNotFoundError = ErrorDetails{
		InternalCode: "notebook:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "notebook not found",
	}
)
