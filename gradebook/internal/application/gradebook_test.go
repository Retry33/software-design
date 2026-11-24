package application_test

import (
	"errors"
	"math"
	"sort"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/application"
	"github.com/Retry33/software-design/gradebook/internal/domain"
)

type fakeRubrics struct{ rubric *domain.Rubric }

func (f *fakeRubrics) Get() (*domain.Rubric, error) {
	if f.rubric == nil {
		f.rubric = domain.NewRubric()
	}
	return f.rubric, nil
}

func (f *fakeRubrics) Save(r *domain.Rubric) error { f.rubric = r; return nil }

type fakeStudents struct {
	byID map[domain.StudentID]*domain.Student
}

func (f *fakeStudents) ByID(id domain.StudentID) (*domain.Student, error) {
	s, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrStudentNotFound
	}
	return s, nil
}

func (f *fakeStudents) All() ([]*domain.Student, error) {
	all := make([]*domain.Student, 0, len(f.byID))
	for _, s := range f.byID {
		all = append(all, s)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID().Value() < all[j].ID().Value() })
	return all, nil
}

func (f *fakeStudents) Save(s *domain.Student) error { f.byID[s.ID()] = s; return nil }

type spyExporter struct {
	destination string
	totals      []application.TotalView
}

func (s *spyExporter) Export(destination string, totals []application.TotalView) error {
	s.destination, s.totals = destination, totals
	return nil
}

func newGradebook(t *testing.T) (*application.Gradebook, *spyExporter) {
	t.Helper()
	students := &fakeStudents{byID: map[domain.StudentID]*domain.Student{}}
	rubrics := &fakeRubrics{}
	exporter := &spyExporter{}
	calculator := domain.NewTotalCalculator(domain.NewRoundedPolicy(domain.StrictPolicy{}, 2))

	return application.NewGradebook(
		application.NewStudentService(students),
		application.NewRubricService(rubrics),
		application.NewGradingService(students, rubrics),
		application.NewReportService(students, rubrics, calculator, exporter),
	), exporter
}

func seed(t *testing.T, g *application.Gradebook) {
	t.Helper()
	for _, cmd := range []application.AddStudentCommand{{ID: "s1", Name: "Ivan Petrov"}, {ID: "s2", Name: "Anna Ivanova"}} {
		if _, err := g.AddStudent(cmd); err != nil {
			t.Fatalf("AddStudent(%s): %v", cmd.ID, err)
		}
	}
	for _, cmd := range []application.AddCriterionCommand{
		{ID: "hw", Name: "Homework", MaxScore: 10, Weight: 30},
		{ID: "exam", Name: "Exam", MaxScore: 20, Weight: 70},
	} {
		if _, err := g.AddCriterion(cmd); err != nil {
			t.Fatalf("AddCriterion(%s): %v", cmd.ID, err)
		}
	}
}

func mustGrade(t *testing.T, g *application.Gradebook, student, criterion string, score int) {
	t.Helper()
	if err := g.Grade(application.GradeCommand{StudentID: student, CriterionID: criterion, Score: score}); err != nil {
		t.Fatalf("Grade(%s, %s, %d): %v", student, criterion, score, err)
	}
}

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestEndToEndScenario(t *testing.T) {
	// Arrange
	g, exporter := newGradebook(t)
	seed(t, g)
	mustGrade(t, g, "s1", "hw", 8)
	mustGrade(t, g, "s1", "exam", 15)
	mustGrade(t, g, "s2", "hw", 10)

	// Act
	err := g.ExportTotals("totals.csv")

	// Assert
	if err != nil {
		t.Fatalf("ExportTotals: %v", err)
	}
	if exporter.destination != "totals.csv" {
		t.Errorf("destination = %q, want totals.csv", exporter.destination)
	}
	want := []application.TotalView{
		{StudentID: "s1", StudentName: "Ivan Petrov", Total: 76.5},
		{StudentID: "s2", StudentName: "Anna Ivanova", Total: 30},
	}
	if len(exporter.totals) != len(want) {
		t.Fatalf("exported %d totals, want %d", len(exporter.totals), len(want))
	}
	for i := range want {
		got := exporter.totals[i]
		if got.StudentID != want[i].StudentID || got.StudentName != want[i].StudentName ||
			!almostEqual(got.Total, want[i].Total) {
			t.Errorf("totals[%d] = %+v, want %+v", i, got, want[i])
		}
	}
}

func TestAddStudentRejectsDuplicateID(t *testing.T) {
	// Arrange
	g, _ := newGradebook(t)
	seed(t, g)

	// Act
	_, err := g.AddStudent(application.AddStudentCommand{ID: "s1", Name: "Someone Else"})

	// Assert
	if !errors.Is(err, domain.ErrStudentAlreadyExists) {
		t.Fatalf("err = %v, want %v", err, domain.ErrStudentAlreadyExists)
	}
}

func TestAddCriterionRejectsInvalidInputAndKeepsRubric(t *testing.T) {
	tests := []struct {
		name    string
		cmd     application.AddCriterionCommand
		wantErr error
	}{
		{"rejects weights above 100", application.AddCriterionCommand{ID: "bonus", Name: "Bonus", MaxScore: 5, Weight: 1}, domain.ErrWeightsExceedLimit},
		{"rejects duplicate id", application.AddCriterionCommand{ID: "hw", Name: "Homework 2", MaxScore: 5, Weight: 1}, domain.ErrCriterionAlreadyExists},
		{"rejects zero max score", application.AddCriterionCommand{ID: "quiz", Name: "Quiz", MaxScore: 0, Weight: 1}, domain.ErrInvalidMaxScore},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			g, _ := newGradebook(t)
			seed(t, g)

			// Act
			_, err := g.AddCriterion(tt.cmd)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			criteria, err := g.Criteria()
			if err != nil {
				t.Fatalf("Criteria: %v", err)
			}
			if len(criteria) != 2 {
				t.Errorf("len(Criteria()) = %d, want 2", len(criteria))
			}
		})
	}
}

func TestGradeRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name    string
		cmd     application.GradeCommand
		wantErr error
	}{
		{"rejects unknown student", application.GradeCommand{StudentID: "nobody", CriterionID: "hw", Score: 5}, domain.ErrStudentNotFound},
		{"rejects criterion outside rubric", application.GradeCommand{StudentID: "s1", CriterionID: "missing", Score: 5}, domain.ErrCriterionNotFound},
		{"rejects score above criterion max", application.GradeCommand{StudentID: "s1", CriterionID: "hw", Score: 12}, domain.ErrScoreOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			g, _ := newGradebook(t)
			seed(t, g)

			// Act
			err := g.Grade(tt.cmd)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			evaluations, err := g.Evaluations("s1")
			if err != nil {
				t.Fatalf("Evaluations: %v", err)
			}
			if len(evaluations) != 0 {
				t.Errorf("evaluations stored despite error: %+v", evaluations)
			}
		})
	}
}

func TestCorrectAndCancelChangeTheTotal(t *testing.T) {
	// Arrange
	g, _ := newGradebook(t)
	seed(t, g)
	mustGrade(t, g, "s1", "hw", 8)
	total := func() float64 {
		t.Helper()
		v, err := g.Total("s1")
		if err != nil {
			t.Fatalf("Total: %v", err)
		}
		return v.Total
	}

	// Act + Assert
	if err := g.Correct(application.GradeCommand{StudentID: "s1", CriterionID: "hw", Score: 10}); err != nil {
		t.Fatalf("Correct: %v", err)
	}
	if got := total(); !almostEqual(got, 30) {
		t.Errorf("after Correct: total = %v, want 30", got)
	}

	if err := g.Cancel(application.CancelEvaluationCommand{StudentID: "s1", CriterionID: "hw"}); err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if got := total(); !almostEqual(got, 0) {
		t.Errorf("after Cancel: total = %v, want 0", got)
	}
}

func TestEvaluationsIncludeCriterionDetails(t *testing.T) {
	// Arrange
	g, _ := newGradebook(t)
	seed(t, g)
	mustGrade(t, g, "s1", "hw", 8)

	// Act
	got, err := g.Evaluations("s1")

	// Assert
	if err != nil {
		t.Fatalf("Evaluations: %v", err)
	}
	want := application.EvaluationView{CriterionID: "hw", CriterionName: "Homework", Score: 8, MaxScore: 10}
	if len(got) != 1 || got[0] != want {
		t.Errorf("Evaluations() = %+v, want [%+v]", got, want)
	}
}
