package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestNewCriterionID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"rejects empty string", "", "", domain.ErrInvalidCriterionID},
		{"rejects whitespace only", "   ", "", domain.ErrInvalidCriterionID},
		{"accepts regular value", "hw", "hw", nil},
		{"trims surrounding whitespace", "  hw  ", "hw", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			id, err := domain.NewCriterionID(tt.input)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && id.Value() != tt.want {
				t.Errorf("Value() = %q, want %q", id.Value(), tt.want)
			}
		})
	}
}

func TestCriterionIDEqualityByValue(t *testing.T) {
	// Arrange
	a, errA := domain.NewCriterionID("hw")
	b, errB := domain.NewCriterionID(" hw ")
	if errA != nil || errB != nil {
		t.Fatalf("unexpected errors: %v, %v", errA, errB)
	}

	// Assert
	if a != b {
		t.Errorf("ids with equal value are not equal: %v != %v", a, b)
	}
}
