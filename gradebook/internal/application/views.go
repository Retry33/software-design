package application

import "github.com/Retry33/software-design/gradebook/internal/domain"

type StudentView struct {
	ID   string
	Name string
}

type CriterionView struct {
	ID       string
	Name     string
	MaxScore int
	Weight   int
}

type EvaluationView struct {
	CriterionID   string
	CriterionName string
	Score         int
	MaxScore      int
}

type TotalView struct {
	StudentID   string
	StudentName string
	Total       float64
}

func toStudentView(s *domain.Student) StudentView {
	return StudentView{ID: s.ID().Value(), Name: s.Name()}
}

func toCriterionView(c *domain.Criterion) CriterionView {
	return CriterionView{
		ID:       c.ID().Value(),
		Name:     c.Name(),
		MaxScore: c.MaxScore().Value(),
		Weight:   c.Weight().Value(),
	}
}
