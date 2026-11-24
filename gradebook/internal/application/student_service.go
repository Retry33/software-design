package application

import (
	"errors"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

type StudentService struct {
	students domain.StudentRepository
}

func NewStudentService(students domain.StudentRepository) *StudentService {
	return &StudentService{students: students}
}

func (s *StudentService) AddStudent(cmd AddStudentCommand) (StudentView, error) {
	id, err := domain.NewStudentID(cmd.ID)
	if err != nil {
		return StudentView{}, err
	}

	_, err = s.students.ByID(id)
	switch {
	case err == nil:
		return StudentView{}, domain.ErrStudentAlreadyExists
	case !errors.Is(err, domain.ErrStudentNotFound):
		return StudentView{}, err
	}

	student, err := domain.NewStudent(id, cmd.Name)
	if err != nil {
		return StudentView{}, err
	}
	if err := s.students.Save(student); err != nil {
		return StudentView{}, err
	}

	return toStudentView(student), nil
}

func (s *StudentService) Students() ([]StudentView, error) {
	students, err := s.students.All()
	if err != nil {
		return nil, err
	}

	views := make([]StudentView, 0, len(students))
	for _, student := range students {
		views = append(views, toStudentView(student))
	}

	return views, nil
}

func loadStudent(students domain.StudentRepository, rawID string) (*domain.Student, error) {
	id, err := domain.NewStudentID(rawID)
	if err != nil {
		return nil, err
	}

	return students.ByID(id)
}
