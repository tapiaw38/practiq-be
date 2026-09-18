package courselevel

import (
	"testing"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func titles(sheets []domain.PracticeSheet) []string {
	out := make([]string, 0, len(sheets))
	for _, sheet := range sheets {
		out = append(out, sheet.Title)
	}
	return out
}

func TestSortSheetsForPath(t *testing.T) {
	base := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)

	t.Run("oldest first", func(t *testing.T) {
		// The listing arrives newest-first, which is what the teacher sees.
		sheets := []domain.PracticeSheet{
			{ID: "c", Title: "Practica 3", CreatedAt: base.Add(2 * time.Hour)},
			{ID: "b", Title: "Practica 2", CreatedAt: base.Add(time.Hour)},
			{ID: "a", Title: "Practica 1", CreatedAt: base},
		}
		sortSheetsForPath(sheets)

		want := []string{"Practica 1", "Practica 2", "Practica 3"}
		for i, title := range want {
			if sheets[i].Title != title {
				t.Fatalf("got %v, want %v", titles(sheets), want)
			}
		}
	})

	t.Run("same instant falls back to id", func(t *testing.T) {
		// Whatever order the database hands these back in, the path is the same.
		forward := []domain.PracticeSheet{
			{ID: "a1", Title: "A", CreatedAt: base},
			{ID: "b2", Title: "B", CreatedAt: base},
		}
		backward := []domain.PracticeSheet{
			{ID: "b2", Title: "B", CreatedAt: base},
			{ID: "a1", Title: "A", CreatedAt: base},
		}
		sortSheetsForPath(forward)
		sortSheetsForPath(backward)

		if forward[0].ID != "a1" || backward[0].ID != "a1" {
			t.Fatalf("unstable: forward=%v backward=%v", titles(forward), titles(backward))
		}
	})
}

func TestSortNotebooksForPath(t *testing.T) {
	base := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	notebooks := []domain.Notebook{
		{ID: "z", Title: "Cuaderno 2", CreatedAt: base.Add(time.Hour)},
		{ID: "y", Title: "Cuaderno 1", CreatedAt: base},
	}

	sortNotebooksForPath(notebooks)

	if notebooks[0].Title != "Cuaderno 1" {
		t.Fatalf("first = %q, want %q", notebooks[0].Title, "Cuaderno 1")
	}
}
