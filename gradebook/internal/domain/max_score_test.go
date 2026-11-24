package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestNewMaxScore(t *testing.T) {
	tests := []struct {
		name     string
		maxValue int
		wantErr  error
	}{
		{"rejects negative value", -1, domain.ErrInvalidMaxScore},
		{"rejects zero value", 0, domain.ErrInvalidMaxScore},
		{"accepts positive value", 10, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			maxScore, err := domain.NewMaxScore(tt.maxValue)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && maxScore.Value() != tt.maxValue {
				t.Errorf("Value() = %v, want %v", maxScore.Value(), tt.maxValue)
			}
		})
	}
}
