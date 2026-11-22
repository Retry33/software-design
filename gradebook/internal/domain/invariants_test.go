package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestDomainConstructorsRejectUninitializedValues(t *testing.T) {
	if _, err := domain.NewStudent(domain.StudentID{}, "Иван"); !errors.Is(err, domain.ErrInvalidStudentID) {
		t.Fatalf("NewStudent: got %v, want %v", err, domain.ErrInvalidStudentID)
	}
	valid := mustCriterion(t, "hw", 10, 30)
	tests := []struct {
		name   string
		id     domain.CriterionID
		max    domain.MaxScore
		weight domain.Weight
		want   error
	}{
		{"id", domain.CriterionID{}, valid.MaxScore(), valid.Weight(), domain.ErrInvalidCriterionID},
		{"max score", valid.ID(), domain.MaxScore{}, valid.Weight(), domain.ErrInvalidMaxScore},
		{"weight", valid.ID(), valid.MaxScore(), domain.Weight{}, domain.ErrWeightOutOfRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewCriterion(tt.id, "Домашние задания", tt.max, tt.weight)
			if !errors.Is(err, tt.want) {
				t.Fatalf("NewCriterion: got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestRubricRejectsUninitializedCriterion(t *testing.T) {
	rubric := domain.NewRubric()
	for _, criterion := range []*domain.Criterion{nil, {}} {
		if err := rubric.AddCriterion(criterion); !errors.Is(err, domain.ErrInvalidCriterionID) {
			t.Fatalf("AddCriterion: got %v, want %v", err, domain.ErrInvalidCriterionID)
		}
	}
	if got := len(rubric.Criteria()); got != 0 {
		t.Fatalf("rubric contains %d criteria, want 0", got)
	}
}

func TestFailedRubricBuildLeavesSourceUntouched(t *testing.T) {
	base, err := domain.NewRubricBuilder().Criterion("hw", "Homework", 10, 30).Build()
	if err != nil {
		t.Fatal(err)
	}
	_, err = domain.NewRubricBuilderFrom(base).
		Criterion("test", "Test", 10, 30).
		Criterion("exam", "Exam", 20, 50).
		Build()
	if !errors.Is(err, domain.ErrWeightsExceedLimit) {
		t.Fatalf("Build: got %v, want %v", err, domain.ErrWeightsExceedLimit)
	}
	if got := len(base.Criteria()); got != 1 {
		t.Fatalf("source contains %d criteria, want 1", got)
	}
}
