package teacherstudentassignment

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	OperationResultData struct {
		Message string `json:"message"`
	}

	UserData struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		ProfileType    string `json:"profile_type"`
		AcademicStatus string `json:"academic_status"`
	}
)

func toUserData(user domain.UserProfile) UserData {
	return UserData{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		ProfileType:    user.ProfileType,
		AcademicStatus: user.AcademicStatus,
	}
}

func toOperationResultData(result domain.OperationResult) OperationResultData {
	return OperationResultData{Message: result.Message}
}
