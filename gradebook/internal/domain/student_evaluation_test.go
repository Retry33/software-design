package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestStudentGrade(t *testing.T) {
	tests := []struct {
		name        string
		criterionID string
		value       int
		wantErr     error
	}{
		{"accepts score inside range", "hw", 8, nil},
		{"accepts score equal to criterion max", "hw", 10, nil},
		{"uses the max of the graded criterion", "exam", 15, nil},
		{"rejects criterion outside rubric", "missing", 5, domain.ErrCriterionNotFound},
		{"rejects score above criterion max", "hw", 11, domain.ErrScoreOutOfRange},
		{"rejects negative score", "hw", -1, domain.ErrScoreOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rubric := newTestRubric(t)
			student := mustStudent(t)

			// Act
			err := student.Grade(rubric, mustCriterionID(t, tt.criterionID), tt.value)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			got := student.Evaluations()
			if err != nil {
				if len(got) != 0 {
					t.Fatalf("evaluation stored despite error: %v", got)
				}
				return
			}
			if len(got) != 1 || got[0].Score().Value() != tt.value {
				t.Errorf("Evaluations() = %v, want one evaluation with score %d", got, tt.value)
			}
		})
	}
}

func TestStudentGradeRejectsSecondEvaluationForSameCriterion(t *testing.T) {
	// Arrange
	rubric := newTestRubric(t)
	student := mustStudent(t)
	hw := mustCriterionID(t, "hw")
	if err := student.Grade(rubric, hw, 8); err != nil {
		t.Fatalf("Grade: %v", err)
	}

	// Act
	err := student.Grade(rubric, hw, 9)

	// Assert
	if !errors.Is(err, domain.ErrEvaluationExists) {
		t.Fatalf("err = %v, want %v", err, domain.ErrEvaluationExists)
	}
	if got := student.Evaluations()[0].Score().Value(); got != 8 {
		t.Errorf("score = %d, want the original 8", got)
	}
}

func TestStudentCorrect(t *testing.T) {
	tests := []struct {
		name      string
		graded    bool
		value     int
		wantErr   error
		wantScore int
	}{
		{"changes existing score", true, 9, nil, 9},
		{"rejects missing evaluation", false, 9, domain.ErrEvaluationNotFound, 0},
		{"rejects out of range and keeps old score", true, 11, domain.ErrScoreOutOfRange, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rubric := newTestRubric(t)
			student := mustStudent(t)
			hw := mustCriterionID(t, "hw")
			if tt.graded {
				if err := student.Grade(rubric, hw, 8); err != nil {
					t.Fatalf("Grade: %v", err)
				}
			}

			// Act
			err := student.Correct(rubric, hw, tt.value)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.graded {
				if got := student.Evaluations()[0].Score().Value(); got != tt.wantScore {
					t.Errorf("score = %d, want %d", got, tt.wantScore)
				}
			}
		})
	}
}

func TestStudentCancel(t *testing.T) {
	// Arrange
	rubric := newTestRubric(t)
	student := mustStudent(t)
	hw := mustCriterionID(t, "hw")
	if err := student.Grade(rubric, hw, 8); err != nil {
		t.Fatalf("Grade: %v", err)
	}

	// Act
	err := student.Cancel(hw)

	// Assert
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if got := len(student.Evaluations()); got != 0 {
		t.Errorf("len(Evaluations()) = %d, want 0", got)
	}
	if err := student.Grade(rubric, hw, 9); err != nil {
		t.Errorf("Grade after Cancel: %v", err)
	}
}

func TestStudentCancelRejectsMissingEvaluation(t *testing.T) {
	// Arrange
	student := mustStudent(t)

	// Act
	err := student.Cancel(mustCriterionID(t, "hw"))

	// Assert
	if !errors.Is(err, domain.ErrEvaluationNotFound) {
		t.Fatalf("err = %v, want %v", err, domain.ErrEvaluationNotFound)
	}
}

func TestStudentEvaluationsAreOrderedByCriterionID(t *testing.T) {
	// Arrange
	rubric := newTestRubric(t)
	student := mustStudent(t)
	if err := student.Grade(rubric, mustCriterionID(t, "hw"), 8); err != nil {
		t.Fatalf("Grade: %v", err)
	}
	if err := student.Grade(rubric, mustCriterionID(t, "exam"), 15); err != nil {
		t.Fatalf("Grade: %v", err)
	}

	// Act
	got := student.Evaluations()

	// Assert
	if len(got) != 2 || got[0].CriterionID().Value() != "exam" || got[1].CriterionID().Value() != "hw" {
		t.Errorf("Evaluations() = %v, want exam then hw", got)
	}
}
