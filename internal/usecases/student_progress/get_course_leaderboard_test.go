package studentprogress

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
)

func TestShortDisplayNameKeepsFirstNameAndInitialOnly(t *testing.T) {
	cases := []struct {
		first, last, want string
	}{
		{"Walter", "Tapia", "Walter T."},
		{"  Sofia ", " Martinez ", "Sofia M."},

		{"Walter", "", "Walter"},

		{"", "Tapia", "Tapia"},

		{"", "", "Alumno"},
		{"   ", "  ", "Alumno"},
	}

	for _, tc := range cases {
		got := shortDisplayName(authapi.UserInfo{FirstName: tc.first, LastName: tc.last})
		if got != tc.want {
			t.Errorf("shortDisplayName(%q, %q) = %q, esperaba %q", tc.first, tc.last, got, tc.want)
		}
	}
}

func TestShortDisplayNameUppercasesTheInitial(t *testing.T) {
	if got := shortDisplayName(authapi.UserInfo{FirstName: "juan", LastName: "perez"}); got != "juan P." {
		t.Fatalf("la inicial del apellido va en mayuscula, dio %q", got)
	}
}

func TestShortDisplayNameHandlesMultibyteSurnames(t *testing.T) {

	if got := shortDisplayName(authapi.UserInfo{FirstName: "Ana", LastName: "Ñandú"}); got != "Ana Ñ." {
		t.Fatalf("una inicial multibyte debe salir entera, dio %q", got)
	}
}
