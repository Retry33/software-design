package application

import "github.com/Retry33/software-design/gradebook/internal/domain"

type RubricService struct {
	rubrics domain.RubricRepository
}

func NewRubricService(rubrics domain.RubricRepository) *RubricService {
	return &RubricService{rubrics: rubrics}
}

func (s *RubricService) AddCriterion(cmd AddCriterionCommand) (CriterionView, error) {
	rubric, err := s.rubrics.Get()
	if err != nil {
		return CriterionView{}, err
	}

	rubric, err = domain.NewRubricBuilderFrom(rubric).
		Criterion(cmd.ID, cmd.Name, cmd.MaxScore, cmd.Weight).
		Build()
	if err != nil {
		return CriterionView{}, err
	}
	if err := s.rubrics.Save(rubric); err != nil {
		return CriterionView{}, err
	}

	criteria := rubric.Criteria()

	return toCriterionView(criteria[len(criteria)-1]), nil
}

func (s *RubricService) Criteria() ([]CriterionView, error) {
	rubric, err := s.rubrics.Get()
	if err != nil {
		return nil, err
	}

	criteria := rubric.Criteria()
	views := make([]CriterionView, 0, len(criteria))
	for _, c := range criteria {
		views = append(views, toCriterionView(c))
	}

	return views, nil
}
