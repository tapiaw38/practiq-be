package mappings

import "net/http"

var (
	RequestBodyParsingError = ErrorDetails{
		Code:       "common:request:body-parsing-error",
		StatusCode: http.StatusBadRequest,
		Message:    "invalid request body",
	}

	InternalServerError = ErrorDetails{
		Code:       "common:internal-server-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}

	UnauthorizedError = ErrorDetails{
		Code:       "common:unauthorized",
		StatusCode: http.StatusUnauthorized,
		Message:    "unauthorized",
	}

	ForbiddenError = ErrorDetails{
		Code:       "common:forbidden",
		StatusCode: http.StatusForbidden,
		Message:    "forbidden",
	}

	NotFoundError = ErrorDetails{
		Code:       "common:not-found",
		StatusCode: http.StatusNotFound,
		Message:    "resource not found",
	}

	// Learning Strategy errors
	LearningStrategyListError = ErrorDetails{
		Code:       "learning-strategy:list-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list learning strategies",
	}

	LearningStrategyGetError = ErrorDetails{
		Code:       "learning-strategy:get-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get learning strategy",
	}

	LearningStrategyCreateError = ErrorDetails{
		Code:       "learning-strategy:create-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to create learning strategy",
	}

	LearningStrategyUpdateError = ErrorDetails{
		Code:       "learning-strategy:update-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to update learning strategy",
	}

	LearningStrategyDeleteError = ErrorDetails{
		Code:       "learning-strategy:delete-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to delete learning strategy",
	}

	LearningStrategyAssignError = ErrorDetails{
		Code:       "learning-strategy:assign-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to assign learning strategy to course",
	}

	LearningStrategyUnassignError = ErrorDetails{
		Code:       "learning-strategy:unassign-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to unassign learning strategy from course",
	}

	// Course Progress errors
	CourseProgressGetError = ErrorDetails{
		Code:       "course-progress:get-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get course progress",
	}

	CourseProgressListError = ErrorDetails{
		Code:       "course-progress:list-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to list course progress",
	}

	// Notebook Submission errors
	NotebookSubmissionGetError = ErrorDetails{
		Code:       "notebook-submission:get-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to get notebook submission",
	}

	NotebookSubmissionReviewError = ErrorDetails{
		Code:       "notebook-submission:review-error",
		StatusCode: http.StatusInternalServerError,
		Message:    "failed to review notebook submission",
	}
)
