package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestRubricAddCriterion(t *testing.T) {
	tests := []struct {
		name            string
		existingWeights []int
		newID           string
		newWeight       int
		wantErr         error
	}{
		{"accepts first criterion", nil, "hw", 30, nil},
		{"accepts weights summing to 100", []int{70}, "new", 30, nil},
		{"rejects weights above 100", []int{70}, "new", 31, domain.ErrWeightsExceedLimit},
		{"rejects duplicate criterion id", []int{30}, "c0", 10, domain.ErrCriterionAlreadyExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rubric := domain.NewRubric()
			for i, w := range tt.existingWeights {
				if err := rubric.AddCriterion(mustCriterion(t, fmt.Sprintf("c%d", i), 10, w)); err != nil {
					t.Fatalf("AddCriterion: %v", err)
				}
			}

			// Act
			err := rubric.AddCriterion(mustCriterion(t, tt.newID, 10, tt.newWeight))

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			want := len(tt.existingWeights)
			if err == nil {
				want++
			}
			if got := len(rubric.Criteria()); got != want {
				t.Errorf("len(Criteria()) = %d, want %d", got, want)
			}
		})
	}
}

func TestRubricCriterion(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"finds registered criterion", "hw", nil},
		{"rejects criterion outside rubric", "missing", domain.ErrCriterionNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			rubric := newTestRubric(t)

			// Act
			c, err := rubric.Criterion(mustCriterionID(t, tt.id))

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && c.ID().Value() != tt.id {
				t.Errorf("ID() = %q, want %q", c.ID().Value(), tt.id)
			}
		})
	}
}

func TestRubricCriteriaReturnsCopy(t *testing.T) {
	// Arrange
	rubric := newTestRubric(t)

	// Act
	got := rubric.Criteria()
	got[0] = nil

	// Assert
	if rubric.Criteria()[0] == nil {
		t.Error("changing the returned slice changed the rubric")
	}
}
