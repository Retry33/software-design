package memory

import (
	"sort"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

var (
	_ domain.StudentRepository = (*StudentRepository)(nil)
	_ domain.RubricRepository  = (*RubricRepository)(nil)
)

type StudentRepository struct {
	students map[domain.StudentID]*domain.Student
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{students: make(map[domain.StudentID]*domain.Student)}
}

func (r *StudentRepository) ByID(id domain.StudentID) (*domain.Student, error) {
	s, ok := r.students[id]
	if !ok {
		return nil, domain.ErrStudentNotFound
	}

	return s, nil
}

func (r *StudentRepository) All() ([]*domain.Student, error) {
	all := make([]*domain.Student, 0, len(r.students))
	for _, s := range r.students {
		all = append(all, s)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID().Value() < all[j].ID().Value() })

	return all, nil
}

func (r *StudentRepository) Save(s *domain.Student) error {
	r.students[s.ID()] = s

	return nil
}

type RubricRepository struct {
	rubric *domain.Rubric
}

func NewRubricRepository() *RubricRepository {
	return &RubricRepository{rubric: domain.NewRubric()}
}

func (r *RubricRepository) Get() (*domain.Rubric, error) {
	return r.rubric, nil
}

func (r *RubricRepository) Save(rubric *domain.Rubric) error {
	r.rubric = rubric

	return nil
}
