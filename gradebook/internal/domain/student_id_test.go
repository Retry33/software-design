package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestNewStudentID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"rejects empty string", "", "", domain.ErrInvalidStudentID},
		{"rejects whitespace only", "   ", "", domain.ErrInvalidStudentID},
		{"accepts regular value", "s1", "s1", nil},
		{"trims surrounding whitespace", "  s1  ", "s1", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			id, err := domain.NewStudentID(tt.input)

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

func TestStudentIDEqualityByValue(t *testing.T) {
	// Arrange
	a, errA := domain.NewStudentID("s1")
	b, errB := domain.NewStudentID(" s1 ")
	if errA != nil || errB != nil {
		t.Fatalf("unexpected errors: %v, %v", errA, errB)
	}

	// Assert
	if a != b {
		t.Errorf("ids with equal value are not equal: %v != %v", a, b)
	}
}
