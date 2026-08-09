package mappings

import "net/http"

var (
	RequestBodyParsingError = ErrorDetails{
		InternalCode: "common:request:body-parsing-error",
		StatusCode:   http.StatusBadRequest,
		Message:      "invalid request body",
	}

	InternalServerError = ErrorDetails{
		InternalCode: "common:internal-server-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "internal server error",
	}

	UnauthorizedError = ErrorDetails{
		InternalCode: "common:unauthorized",
		StatusCode:   http.StatusUnauthorized,
		Message:      "unauthorized",
	}

	ForbiddenError = ErrorDetails{
		InternalCode: "common:forbidden",
		StatusCode:   http.StatusForbidden,
		Message:      "forbidden",
	}

	NotFoundError = ErrorDetails{
		InternalCode: "common:not-found",
		StatusCode:   http.StatusNotFound,
		Message:      "resource not found",
	}

	// Learning Strategy errors
	LearningStrategyListError = ErrorDetails{
		InternalCode: "learning-strategy:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list learning strategies",
	}

	LearningStrategyGetError = ErrorDetails{
		InternalCode: "learning-strategy:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get learning strategy",
	}

	LearningStrategyCreateError = ErrorDetails{
		InternalCode: "learning-strategy:create-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to create learning strategy",
	}

	LearningStrategyUpdateError = ErrorDetails{
		InternalCode: "learning-strategy:update-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to update learning strategy",
	}

	LearningStrategyDeleteError = ErrorDetails{
		InternalCode: "learning-strategy:delete-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to delete learning strategy",
	}

	LearningStrategyAssignError = ErrorDetails{
		InternalCode: "learning-strategy:assign-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to assign learning strategy to course",
	}

	LearningStrategyUnassignError = ErrorDetails{
		InternalCode: "learning-strategy:unassign-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to unassign learning strategy from course",
	}

	// Course Progress errors
	CourseProgressGetError = ErrorDetails{
		InternalCode: "course-progress:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get course progress",
	}

	CourseProgressListError = ErrorDetails{
		InternalCode: "course-progress:list-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to list course progress",
	}

	// Notebook Submission errors
	NotebookSubmissionGetError = ErrorDetails{
		InternalCode: "notebook-submission:get-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to get notebook submission",
	}

	NotebookSubmissionReviewError = ErrorDetails{
		InternalCode: "notebook-submission:review-error",
		StatusCode:   http.StatusInternalServerError,
		Message:      "failed to review notebook submission",
	}
)
