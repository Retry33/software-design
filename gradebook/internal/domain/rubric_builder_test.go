package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

type rawCriterion struct {
	id, name         string
	maxScore, weight int
}

func TestRubricBuilder(t *testing.T) {
	tests := []struct {
		name         string
		criteria     []rawCriterion
		wantErr      error
		wantCriteria int
	}{
		{"builds rubric from raw values",
			[]rawCriterion{{"hw", "Homework", 10, 30}, {"exam", "Exam", 20, 70}}, nil, 2},
		{"stops at invalid max score",
			[]rawCriterion{{"hw", "Homework", 0, 30}}, domain.ErrInvalidMaxScore, 0},
		{"stops at empty criterion name",
			[]rawCriterion{{"hw", "  ", 10, 30}}, domain.ErrInvalidCriterionName, 0},
		{"stops at weights above 100",
			[]rawCriterion{{"hw", "Homework", 10, 60}, {"exam", "Exam", 20, 50}}, domain.ErrWeightsExceedLimit, 0},
		{"keeps the first error",
			[]rawCriterion{{"", "Homework", 10, 30}, {"exam", "", 20, 70}}, domain.ErrInvalidCriterionID, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			builder := domain.NewRubricBuilder()
			for _, c := range tt.criteria {
				builder.Criterion(c.id, c.name, c.maxScore, c.weight)
			}

			// Act
			rubric, err := builder.Build()

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				if rubric != nil {
					t.Error("Build() returned a rubric together with an error")
				}
				return
			}
			if got := len(rubric.Criteria()); got != tt.wantCriteria {
				t.Errorf("len(Criteria()) = %d, want %d", got, tt.wantCriteria)
			}
		})
	}
}
