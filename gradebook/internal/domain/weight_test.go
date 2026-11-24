package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func TestNewWeight(t *testing.T) {
	tests := []struct {
		name    string
		weight  int
		wantErr error
	}{
		{"rejects negative value", -1, domain.ErrWeightOutOfRange},
		{"rejects zero value", 0, domain.ErrWeightOutOfRange},
		{"accepts value inside range", 10, nil},
		{"rejects value above 100", 101, domain.ErrWeightOutOfRange},
		{"accepts lower bound", 1, nil},
		{"accepts upper bound", 100, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			weight, err := domain.NewWeight(tt.weight)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && weight.Value() != tt.weight {
				t.Errorf("Value() = %v, want %v", weight.Value(), tt.weight)
			}
		})
	}
}
