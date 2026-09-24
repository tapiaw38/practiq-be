package practicesheet

import (
	"context"
	"errors"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories"
	studentcoursexp "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/student_course_xp"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type fakeCourseXP struct {
	balance    int
	awarded    []studentcoursexp.AwardInput
	seen       map[string]bool
	awardErr   error
	balanceErr error
}

func newFakeCourseXP(balance int) *fakeCourseXP {
	return &fakeCourseXP{balance: balance, seen: map[string]bool{}}
}

func (f *fakeCourseXP) Award(_ context.Context, input studentcoursexp.AwardInput) (studentcoursexp.AwardResult, error) {
	if f.awardErr != nil {
		return studentcoursexp.AwardResult{}, f.awardErr
	}
	f.awarded = append(f.awarded, input)
	if f.seen[input.EventKey] {
		return studentcoursexp.AwardResult{Awarded: false, TotalXP: f.balance}, nil
	}
	f.seen[input.EventKey] = true
	f.balance += input.Points
	return studentcoursexp.AwardResult{Awarded: true, TotalXP: f.balance}, nil
}

func (f *fakeCourseXP) Balance(context.Context, string, string) (int, error) {
	if f.balanceErr != nil {
		return 0, f.balanceErr
	}
	return f.balance, nil
}

func (f *fakeCourseXP) ListByStudent(context.Context, string) ([]studentcoursexp.CourseXP, error) {
	return nil, nil
}

func (f *fakeCourseXP) LeaderboardByCourse(context.Context, string, string, int) ([]studentcoursexp.LeaderboardEntry, error) {
	return nil, nil
}

func (f *fakeCourseXP) pointsFor(eventType string) int {
	total := 0
	for _, input := range f.awarded {
		if input.EventType == eventType {
			total += input.Points
		}
	}
	return total
}

func xpApp(repo *fakeCourseXP) *appcontext.Context {
	return &appcontext.Context{Repositories: &repositories.Repositories{StudentCourseXP: repo}}
}

func xpSheet(sheetType string, exerciseIDs ...string) *domain.PracticeSheet {
	sheet := &domain.PracticeSheet{ID: "sheet-1", CourseID: "course-1", SheetType: sheetType}
	for _, id := range exerciseIDs {
		sheet.Exercises = append(sheet.Exercises, domain.PracticeSheetExercise{Exercise: domain.Exercise{ID: id}})
	}
	return sheet
}

func answered(ids ...string) []AttemptInput {
	attempts := make([]AttemptInput, 0, len(ids))
	for _, id := range ids {
		attempts = append(attempts, AttemptInput{ExerciseID: id, AnswerText: "una respuesta"})
	}
	return attempts
}

func TestAwardPracticeXPReportsExistingBalanceWhenNothingIsEarned(t *testing.T) {
	repo := newFakeCourseXP(500)
	sheet := xpSheet("practice", "e1", "e2")
	blank := []AttemptInput{{ExerciseID: "e1"}, {ExerciseID: "e2"}}

	xp := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet, blank, nil, false)

	if xp.Gained != 0 {
		t.Fatalf("respuestas en blanco no deben sumar XP, sumo %d", xp.Gained)
	}
	if xp.Balance != 500 {
		t.Fatalf("un alumno con 500 XP no puede reportar %d", xp.Balance)
	}
}

func TestAwardPracticeXPReportsExistingBalanceWhenAwardsFail(t *testing.T) {
	repo := newFakeCourseXP(320)
	repo.awardErr = errors.New("la base se cayo")
	sheet := xpSheet("practice", "e1")

	xp := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet,
		answered("e1"), []xpExercise{{ExerciseID: "e1", Correct: true}}, false)

	if xp.Gained != 0 {
		t.Fatalf("un award fallido no otorga XP, otorgo %d", xp.Gained)
	}
	if xp.Balance != 320 {
		t.Fatalf("el saldo previo debe sobrevivir al fallo, dio %d", xp.Balance)
	}
}

func TestAwardPracticeXPPaysAttemptAndCorrectAndCompletion(t *testing.T) {
	repo := newFakeCourseXP(0)
	sheet := xpSheet("practice", "e1", "e2")
	exercises := []xpExercise{{ExerciseID: "e1", Correct: true}, {ExerciseID: "e2"}}

	xp := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet,
		answered("e1", "e2"), exercises, false)

	// dos intentos (+1 c/u), un acierto (+10) y la practica completa (+10)
	const expected = xpAttempt*2 + xpCorrect + xpPracticeComplete
	if xp.Gained != expected {
		t.Fatalf("esperaba %d XP, dio %d", expected, xp.Gained)
	}
	if xp.Balance != expected {
		t.Fatalf("el saldo debe seguir a lo ganado, dio %d", xp.Balance)
	}
	if got := repo.pointsFor("exercise_attempt"); got != xpAttempt*2 {
		t.Fatalf("cada ejercicio intentado suma %d, dio %d", xpAttempt, got)
	}
}

func TestAwardPracticeXPSkipsCompletionWhenAnExerciseIsUntouched(t *testing.T) {
	repo := newFakeCourseXP(0)
	sheet := xpSheet("practice", "e1", "e2")

	xp := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet,
		answered("e1"), []xpExercise{{ExerciseID: "e1", Correct: true}}, false)

	if got := repo.pointsFor("practice_complete"); got != 0 {
		t.Fatalf("una practica a medias no esta completa, pago %d", got)
	}
	if expected := xpAttempt + xpCorrect; xp.Gained != expected {
		t.Fatalf("esperaba %d XP, dio %d", expected, xp.Gained)
	}
}

func TestAwardPracticeXPDoesNotPayCorrectForUngradedExercise(t *testing.T) {
	repo := newFakeCourseXP(0)
	sheet := xpSheet("practice", "e1")

	awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet,
		answered("e1"), []xpExercise{{ExerciseID: "e1", Correct: true, Ungraded: true}}, false)

	if got := repo.pointsFor("exercise_correct"); got != 0 {
		t.Fatalf("un ejercicio sin corregir no paga acierto, pago %d", got)
	}
}

func TestAwardPracticeXPPaysLevelTestOnlyWhenPassed(t *testing.T) {
	sheet := xpSheet(sheetTypeLevelTest, "e1")

	passed := newFakeCourseXP(0)
	awardPracticeXP(context.Background(), xpApp(passed), "student-1", sheet,
		answered("e1"), []xpExercise{{ExerciseID: "e1", Correct: true}}, true)
	if got := passed.pointsFor("level_test_pass"); got != xpLevelTestPass {
		t.Fatalf("una prueba aprobada paga %d, pago %d", xpLevelTestPass, got)
	}

	failed := newFakeCourseXP(0)
	awardPracticeXP(context.Background(), xpApp(failed), "student-1", sheet,
		answered("e1"), []xpExercise{{ExerciseID: "e1", Correct: true}}, false)
	if got := failed.pointsFor("level_test_pass"); got != 0 {
		t.Fatalf("una prueba desaprobada no paga, pago %d", got)
	}
	if got := failed.pointsFor("practice_complete"); got != 0 {
		t.Fatalf("una prueba de nivel no paga practica completa, pago %d", got)
	}
}

func TestAwardPracticeXPIsIdempotentAcrossResubmissions(t *testing.T) {
	repo := newFakeCourseXP(0)
	sheet := xpSheet("practice", "e1")
	exercises := []xpExercise{{ExerciseID: "e1", Correct: true}}

	first := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet, answered("e1"), exercises, false)
	second := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet, answered("e1"), exercises, false)

	if second.Gained != 0 {
		t.Fatalf("reenviar la misma practica no puede farmear XP, sumo %d", second.Gained)
	}
	if second.Balance != first.Gained {
		t.Fatalf("el saldo no debe moverse en el reenvio: %d vs %d", second.Balance, first.Gained)
	}
}

func TestAwardPracticeXPGroupsBreakdownByReasonInEarnedOrder(t *testing.T) {
	repo := newFakeCourseXP(0)
	sheet := xpSheet("practice", "e1", "e2", "e3")
	exercises := []xpExercise{
		{ExerciseID: "e1", Correct: true},
		{ExerciseID: "e2"},
		{ExerciseID: "e3", Correct: true},
	}

	xp := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet,
		answered("e1", "e2", "e3"), exercises, false)

	want := []XPBreakdownEntry{
		{EventType: "exercise_attempt", Points: xpAttempt * 3, Count: 3},
		{EventType: "exercise_correct", Points: xpCorrect * 2, Count: 2},
		{EventType: "practice_complete", Points: xpPracticeComplete, Count: 1},
	}
	if len(xp.Breakdown) != len(want) {
		t.Fatalf("esperaba %d motivos, dio %d: %+v", len(want), len(xp.Breakdown), xp.Breakdown)
	}
	for i, entry := range want {
		if xp.Breakdown[i] != entry {
			t.Fatalf("motivo %d: esperaba %+v, dio %+v", i, entry, xp.Breakdown[i])
		}
	}

	total := 0
	for _, entry := range xp.Breakdown {
		total += entry.Points
	}
	if total != xp.Gained {
		t.Fatalf("el desglose debe sumar lo ganado: %d vs %d", total, xp.Gained)
	}
}

func TestAwardPracticeXPLeavesBreakdownEmptyOnResubmission(t *testing.T) {
	repo := newFakeCourseXP(0)
	sheet := xpSheet("practice", "e1")
	exercises := []xpExercise{{ExerciseID: "e1", Correct: true}}

	awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet, answered("e1"), exercises, false)
	second := awardPracticeXP(context.Background(), xpApp(repo), "student-1", sheet, answered("e1"), exercises, false)

	if len(second.Breakdown) != 0 {
		t.Fatalf("un reenvio no paga nada, no hay burbujas que mostrar: %+v", second.Breakdown)
	}
}
