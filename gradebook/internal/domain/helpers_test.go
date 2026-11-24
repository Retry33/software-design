package domain_test

import (
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func mustCriterionID(t *testing.T, value string) domain.CriterionID {
	t.Helper()
	id, err := domain.NewCriterionID(value)
	if err != nil {
		t.Fatalf("NewCriterionID(%q): %v", value, err)
	}
	return id
}

func mustCriterion(t *testing.T, id string, maxValue, weight int) *domain.Criterion {
	t.Helper()
	maxScore, err := domain.NewMaxScore(maxValue)
	if err != nil {
		t.Fatalf("NewMaxScore(%d): %v", maxValue, err)
	}
	w, err := domain.NewWeight(weight)
	if err != nil {
		t.Fatalf("NewWeight(%d): %v", weight, err)
	}
	c, err := domain.NewCriterion(mustCriterionID(t, id), "Criterion "+id, maxScore, w)
	if err != nil {
		t.Fatalf("NewCriterion: %v", err)
	}
	return c
}

func mustStudent(t *testing.T) *domain.Student {
	t.Helper()
	id, err := domain.NewStudentID("s1")
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}
	s, err := domain.NewStudent(id, "Ivan Petrov")
	if err != nil {
		t.Fatalf("NewStudent: %v", err)
	}
	return s
}

func newTestRubric(t *testing.T) *domain.Rubric {
	t.Helper()
	r := domain.NewRubric()
	for _, c := range []*domain.Criterion{
		mustCriterion(t, "hw", 10, 30),
		mustCriterion(t, "exam", 20, 70),
	} {
		if err := r.AddCriterion(c); err != nil {
			t.Fatalf("AddCriterion: %v", err)
		}
	}
	return r
}
