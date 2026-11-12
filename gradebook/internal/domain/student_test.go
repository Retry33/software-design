package domain_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func newTestStudent(t *testing.T, name string) (*domain.Student, error) {
	t.Helper()
	id, err := domain.NewStudentID("s1")
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}
	return domain.NewStudent(id, name)
}

func TestNewStudent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantErr  error
	}{
		{"rejects empty name", "", "", domain.ErrInvalidStudentName},
		{"rejects whitespace-only name", "   ", "", domain.ErrInvalidStudentName},
		{"accepts regular name", "Ivan Petrov", "Ivan Petrov", nil},
		{"trims surrounding whitespace in name", "  Ivan Petrov  ", "Ivan Petrov", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			s, err := newTestStudent(t, tt.input)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if s.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", s.Name(), tt.wantName)
			}
		})
	}
}

func TestStudentRename(t *testing.T) {
	tests := []struct {
		name     string
		newName  string
		wantName string
		wantErr  error
	}{
		{"renames to new name", "Anna Ivanova", "Anna Ivanova", nil},
		{"trims whitespace on rename", "  Anna Ivanova  ", "Anna Ivanova", nil},
		{"rejects empty name and keeps old one", "  ", "Ivan Petrov", domain.ErrInvalidStudentName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			s, err := newTestStudent(t, "Ivan Petrov")
			if err != nil {
				t.Fatalf("NewStudent: %v", err)
			}

			// Act
			err = s.Rename(tt.newName)

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if s.Name() != tt.wantName {
				t.Errorf("Name() = %q, want %q", s.Name(), tt.wantName)
			}
		})
	}
}
