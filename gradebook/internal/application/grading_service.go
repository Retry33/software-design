package application

import "github.com/Retry33/software-design/gradebook/internal/domain"

type GradingService struct {
	students domain.StudentRepository
	rubrics  domain.RubricRepository
}

func NewGradingService(students domain.StudentRepository, rubrics domain.RubricRepository) *GradingService {
	return &GradingService{students: students, rubrics: rubrics}
}

func (s *GradingService) Grade(cmd GradeCommand) error {
	return s.change(cmd.StudentID, cmd.CriterionID,
		func(st *domain.Student, r *domain.Rubric, id domain.CriterionID) error {
			return st.Grade(r, id, cmd.Score)
		})
}

func (s *GradingService) Correct(cmd GradeCommand) error {
	return s.change(cmd.StudentID, cmd.CriterionID,
		func(st *domain.Student, r *domain.Rubric, id domain.CriterionID) error {
			return st.Correct(r, id, cmd.Score)
		})
}

func (s *GradingService) Cancel(cmd CancelEvaluationCommand) error {
	return s.change(cmd.StudentID, cmd.CriterionID,
		func(st *domain.Student, _ *domain.Rubric, id domain.CriterionID) error {
			return st.Cancel(id)
		})
}

func (s *GradingService) Evaluations(studentID string) ([]EvaluationView, error) {
	student, err := loadStudent(s.students, studentID)
	if err != nil {
		return nil, err
	}
	rubric, err := s.rubrics.Get()
	if err != nil {
		return nil, err
	}

	evaluations := student.Evaluations()
	views := make([]EvaluationView, 0, len(evaluations))
	for _, e := range evaluations {
		view := EvaluationView{CriterionID: e.CriterionID().Value(), Score: e.Score().Value()}
		if c, err := rubric.Criterion(e.CriterionID()); err == nil {
			view.CriterionName = c.Name()
			view.MaxScore = c.MaxScore().Value()
		}
		views = append(views, view)
	}

	return views, nil
}

func (s *GradingService) change(
	studentID, criterionID string,
	op func(*domain.Student, *domain.Rubric, domain.CriterionID) error,
) error {
	student, err := loadStudent(s.students, studentID)
	if err != nil {
		return err
	}
	id, err := domain.NewCriterionID(criterionID)
	if err != nil {
		return err
	}
	rubric, err := s.rubrics.Get()
	if err != nil {
		return err
	}

	if err := op(student, rubric, id); err != nil {
		return err
	}

	return s.students.Save(student)
}
