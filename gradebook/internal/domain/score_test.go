package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestNewScore(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		maxValue int
		wantErr  error
	}{
		{"accepts zero", 0, 10, nil},
		{"accepts value equal to max", 10, 10, nil},
		{"accepts value inside range", 5, 10, nil},
		{"rejects negative score", -1, 10, domain.ErrScoreOutOfRange},
		{"rejects score above max", 11, 10, domain.ErrScoreOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			maxScore, err := domain.NewMaxScore(tt.maxValue)
			if err != nil {
				t.Fatalf("NewMaxScore(%d): %v", tt.maxValue, err)
			}

			// Act
			got, err := domain.NewScore(tt.value, maxScore)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && got.Value() != tt.value {
				t.Errorf("Value() = %v, want %v", got.Value(), tt.value)
			}
		})
	}
}

func TestNewScoreRejectsZeroValueMaxScore(t *testing.T) {
	_, err := domain.NewScore(0, domain.MaxScore{})
	if !errors.Is(err, domain.ErrInvalidMaxScore) {
		t.Fatalf("err = %v, want %v", err, domain.ErrInvalidMaxScore)
	}
}
