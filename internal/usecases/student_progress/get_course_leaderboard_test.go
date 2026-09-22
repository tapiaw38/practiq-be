package studentprogress

import "testing"

func TestShortDisplayNameKeepsFirstNameAndInitialOnly(t *testing.T) {
	cases := map[string]string{
		"Walter Tapia":        "Walter T.",
		"Walter Tapia Gomez":  "Walter T.",
		"  Sofia   Martinez ": "Sofia M.",
		// A single word is all the student gave: there is no surname to hide.
		"Walter": "Walter",
		// A profile with no name still has to render as something.
		"":    "Alumno",
		"   ": "Alumno",
	}

	for name, want := range cases {
		if got := shortDisplayName(name); got != want {
			t.Errorf("shortDisplayName(%q) = %q, esperaba %q", name, got, want)
		}
	}
}

func TestShortDisplayNameUppercasesTheInitial(t *testing.T) {
	if got := shortDisplayName("juan perez"); got != "juan P." {
		t.Fatalf("la inicial del apellido va en mayuscula, dio %q", got)
	}
}

func TestShortDisplayNameHandlesMultibyteSurnames(t *testing.T) {
	// Slicing bytes instead of runes would cut "Ñ" in half and emit garbage.
	if got := shortDisplayName("Ana Ñandú"); got != "Ana Ñ." {
		t.Fatalf("una inicial multibyte debe salir entera, dio %q", got)
	}
}
