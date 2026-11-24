package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func newTestCriterion(t *testing.T, name string) (*domain.Criterion, error) {
	t.Helper()
	id, err := domain.NewCriterionID("hw")
	if err != nil {
		t.Fatalf("NewCriterionID: %v", err)
	}
	maxScore, err := domain.NewMaxScore(10)
	if err != nil {
		t.Fatalf("NewMaxScore: %v", err)
	}
	weight, err := domain.NewWeight(30)
	if err != nil {
		t.Fatalf("NewWeight: %v", err)
	}
	return domain.NewCriterion(id, name, maxScore, weight)
}

func TestNewCriterion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantErr  error
	}{
		{"rejects empty name", "", "", domain.ErrInvalidCriterionName},
		{"rejects whitespace-only name", "   ", "", domain.ErrInvalidCriterionName},
		{"accepts regular name", "Homework", "Homework", nil},
		{"trims surrounding whitespace in name", "  Homework  ", "Homework", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			c, err := newTestCriterion(t, tt.input)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if c.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", c.Name(), tt.wantName)
			}
		})
	}
}

func TestCriterionRename(t *testing.T) {
	tests := []struct {
		name     string
		newName  string
		wantName string
		wantErr  error
	}{
		{"renames to new name", "Exam", "Exam", nil},
		{"trims whitespace on rename", "  Exam  ", "Exam", nil},
		{"rejects empty name and keeps old one", "  ", "Homework", domain.ErrInvalidCriterionName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, err := newTestCriterion(t, "Homework")
			if err != nil {
				t.Fatalf("NewCriterion: %v", err)
			}

			// Act
			err = c.Rename(tt.newName)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if c.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", c.Name(), tt.wantName)
			}
		})
	}
}
