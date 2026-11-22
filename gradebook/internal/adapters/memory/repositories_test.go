package memory_test

import (
	"errors"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/adapters/memory"
	"github.com/Retry33/software-design/gradebook/internal/domain"
)

func mustStudent(t *testing.T, id, name string) *domain.Student {
	t.Helper()
	sid, err := domain.NewStudentID(id)
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}
	s, err := domain.NewStudent(sid, name)
	if err != nil {
		t.Fatalf("NewStudent: %v", err)
	}
	return s
}

func TestStudentRepositoryByIDReturnsNotFound(t *testing.T) {
	// Arrange
	repo := memory.NewStudentRepository()
	id, err := domain.NewStudentID("nobody")
	if err != nil {
		t.Fatalf("NewStudentID: %v", err)
	}

	// Act
	_, err = repo.ByID(id)

	// Assert
	if !errors.Is(err, domain.ErrStudentNotFound) {
		t.Fatalf("err = %v, want %v", err, domain.ErrStudentNotFound)
	}
}

func TestStudentRepositoryAllIsOrderedByID(t *testing.T) {
	// Arrange
	repo := memory.NewStudentRepository()
	for _, s := range []*domain.Student{mustStudent(t, "s2", "B"), mustStudent(t, "s1", "A")} {
		if err := repo.Save(s); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	// Act
	all, err := repo.All()

	// Assert
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 2 || all[0].ID().Value() != "s1" || all[1].ID().Value() != "s2" {
		t.Errorf("All() order is wrong: got %d students", len(all))
	}
}
