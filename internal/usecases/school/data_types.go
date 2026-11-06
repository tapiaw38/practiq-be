package school

type SchoolData struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Billing string `json:"billing"`
	Status  string `json:"status"`

	Role string `json:"role,omitempty"`

	Owner *SchoolOwner `json:"owner,omitempty"`
	Plan  *SchoolPlan  `json:"plan,omitempty"`
}

type SchoolOwner struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SchoolPlan struct {
	Name        string `json:"name"`
	MaxStudents int    `json:"max_students"`
	Active      bool   `json:"active"`
}

type MemberData struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

type ArchiveCourseData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	GradeName   string `json:"grade_name"`
	SubjectName string `json:"subject_name"`
}
