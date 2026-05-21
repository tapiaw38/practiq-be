package mappings

import "net/http"

var (
	// Course errors
	CourseCreateError = ErrorDetails{
		Code:       "course:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create course",
	}
	CourseGetError = ErrorDetails{
		Code:       "course:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get course",
	}
	CourseNotFoundError = ErrorDetails{
		Code:       "course:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "course not found",
	}
	CourseUpdateError = ErrorDetails{
		Code:       "course:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update course",
	}
	CourseDeleteError = ErrorDetails{
		Code:       "course:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete course",
	}
	CourseListError = ErrorDetails{
		Code:       "course:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list courses",
	}

	// Topic errors
	TopicCreateError = ErrorDetails{
		Code:       "topic:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create topic",
	}
	TopicListError = ErrorDetails{
		Code:       "topic:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list topics",
	}

	// Exercise errors
	ExerciseCreateError = ErrorDetails{
		Code:       "exercise:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create exercise",
	}
	ExerciseListError = ErrorDetails{
		Code:       "exercise:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list exercises",
	}
	ExerciseUpdateError = ErrorDetails{
		Code:       "exercise:update:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update exercise",
	}
	ExerciseDeleteError = ErrorDetails{
		Code:       "exercise:delete:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete exercise",
	}
	ExerciseNotFoundError = ErrorDetails{
		Code:       "exercise:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "exercise not found",
	}

	// Material errors
	MaterialCreateError = ErrorDetails{
		Code:       "material:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create material",
	}
	MaterialListError = ErrorDetails{
		Code:       "material:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list materials",
	}

	// Practice sheet errors
	PracticeSheetCreateError = ErrorDetails{
		Code:       "practice-sheet:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create practice sheet",
	}
	PracticeSheetListError = ErrorDetails{
		Code:       "practice-sheet:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list practice sheets",
	}
	PracticeSheetGetError = ErrorDetails{
		Code:       "practice-sheet:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get practice sheet",
	}
	PracticeSheetNotFoundError = ErrorDetails{
		Code:       "practice-sheet:get:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "practice sheet not found",
	}
	PracticeSheetSubmitError = ErrorDetails{
		Code:       "practice-sheet:submit:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to submit practice sheet",
	}

	// Enrollment errors
	EnrollmentCreateError = ErrorDetails{
		Code:       "enrollment:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to enroll student",
	}
	EnrollmentListError = ErrorDetails{
		Code:       "enrollment:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list students",
	}
	EnrollmentAlreadyExistsError = ErrorDetails{
		Code:       "enrollment:create:already-exists",
		StatusCode: http.StatusConflict,
		Message:    "student already enrolled in this course",
	}

	// Progress errors
	ProgressGetError = ErrorDetails{
		Code:       "progress:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get progress",
	}

	// AI errors
	AIConversationCreateError = ErrorDetails{
		Code:       "ai:conversation:create:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create conversation",
	}
	AIMessageListError = ErrorDetails{
		Code:       "ai:messages:list:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list messages",
	}
	AIHelpError = ErrorDetails{
		Code:       "ai:help:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get AI help",
	}

	// Profile errors
	ProfileSyncError = ErrorDetails{
		Code:       "profile:sync:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to sync profile",
	}
	ProfileGetError = ErrorDetails{
		Code:       "profile:get:error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get profile",
	}
)
