package domain_test

import (
	"math"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestTotalCalculator(t *testing.T) {
	tests := []struct {
		name   string
		policy domain.TotalPolicy
		grades map[string]int
		want   float64
	}{
		{"strict: all criteria graded", domain.StrictPolicy{}, map[string]int{"hw": 8, "exam": 15}, 76.5},
		{"strict: missing evaluation counts as zero", domain.StrictPolicy{}, map[string]int{"hw": 8}, 24},
		{"strict: no evaluations", domain.StrictPolicy{}, nil, 0},
		{"graded only: all criteria graded", domain.GradedOnlyPolicy{}, map[string]int{"hw": 8, "exam": 15}, 76.5},
		{"graded only: normalizes by graded weights", domain.GradedOnlyPolicy{}, map[string]int{"hw": 8}, 80},
		{"graded only: no evaluations", domain.GradedOnlyPolicy{}, nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rubric := newTestRubric(t)
			student := mustStudent(t)
			for id, value := range tt.grades {
				if err := student.Grade(rubric, mustCriterionID(t, id), value); err != nil {
					t.Fatalf("Grade(%s, %d): %v", id, value, err)
				}
			}

			// Act
			got := domain.NewTotalCalculator(tt.policy).Total(rubric, student).Value()

			// Assert
			if !almostEqual(got, tt.want) {
				t.Errorf("Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTotalFollowsCurrentEvaluations(t *testing.T) {
	// Arrange
	rubric := newTestRubric(t)
	student := mustStudent(t)
	hw := mustCriterionID(t, "hw")
	calculator := domain.NewTotalCalculator(domain.StrictPolicy{})
	total := func() float64 { return calculator.Total(rubric, student).Value() }

	if err := student.Grade(rubric, hw, 8); err != nil {
		t.Fatalf("Grade: %v", err)
	}
	if got := total(); !almostEqual(got, 24) {
		t.Errorf("after Grade: total = %v, want 24", got)
	}

	if err := student.Correct(rubric, hw, 10); err != nil {
		t.Fatalf("Correct: %v", err)
	}
	if got := total(); !almostEqual(got, 30) {
		t.Errorf("after Correct: total = %v, want 30", got)
	}

	if err := student.Cancel(hw); err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if got := total(); !almostEqual(got, 0) {
		t.Errorf("after Cancel: total = %v, want 0", got)
	}
}

func TestRoundedPolicy(t *testing.T) {
	rubric := domain.NewRubric()
	if err := rubric.AddCriterion(mustCriterion(t, "quiz", 3, 100)); err != nil {
		t.Fatalf("AddCriterion: %v", err)
	}
	student := mustStudent(t)
	if err := student.Grade(rubric, mustCriterionID(t, "quiz"), 1); err != nil {
		t.Fatalf("Grade: %v", err)
	}

	tests := []struct {
		name   string
		policy domain.TotalPolicy
		want   float64
	}{
		{"keeps unrounded value without decorator", domain.StrictPolicy{}, 100.0 / 3},
		{"rounds to two places", domain.NewRoundedPolicy(domain.StrictPolicy{}, 2), 33.33},
		{"rounds to whole number", domain.NewRoundedPolicy(domain.StrictPolicy{}, 0), 33},
		{"decorates any policy", domain.NewRoundedPolicy(domain.GradedOnlyPolicy{}, 1), 33.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got := domain.NewTotalCalculator(tt.policy).Total(rubric, student).Value()

			// Assert
			if !almostEqual(got, tt.want) {
				t.Errorf("Total() = %v, want %v", got, tt.want)
			}
		})
	}
}
