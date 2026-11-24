package domain

import "sort"

type Student struct {
	id          StudentID
	name        string
	evaluations map[CriterionID]*Evaluation
}

func NewStudent(id StudentID, name string) (*Student, error) {
	if id.Value() == "" {
		return nil, ErrInvalidStudentID
	}
	normalizedName, ok := normalizeName(name)
	if !ok {
		return nil, ErrInvalidStudentName
	}

	return &Student{
		id:          id,
		name:        normalizedName,
		evaluations: make(map[CriterionID]*Evaluation),
	}, nil
}

func (s *Student) Rename(name string) error {
	normalizedName, ok := normalizeName(name)
	if !ok {
		return ErrInvalidStudentName
	}

	s.name = normalizedName

	return nil
}

func (s *Student) ID() StudentID {
	return s.id
}

func (s *Student) Name() string {
	return s.name
}

func (s *Student) Grade(rubric *Rubric, id CriterionID, value int) error {
	score, err := scoreFor(rubric, id, value)
	if err != nil {
		return err
	}
	if _, exists := s.evaluations[id]; exists {
		return ErrEvaluationExists
	}

	s.evaluations[id] = newEvaluation(id, score)

	return nil
}

func (s *Student) Correct(rubric *Rubric, id CriterionID, value int) error {
	evaluation, exists := s.evaluations[id]
	if !exists {
		return ErrEvaluationNotFound
	}

	score, err := scoreFor(rubric, id, value)
	if err != nil {
		return err
	}

	evaluation.score = score

	return nil
}

func (s *Student) Cancel(id CriterionID) error {
	if _, exists := s.evaluations[id]; !exists {
		return ErrEvaluationNotFound
	}

	delete(s.evaluations, id)

	return nil
}

func (s *Student) Evaluations() []Evaluation {
	result := make([]Evaluation, 0, len(s.evaluations))
	for _, e := range s.evaluations {
		result = append(result, *e)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].criterionID.Value() < result[j].criterionID.Value()
	})

	return result
}

func scoreFor(rubric *Rubric, id CriterionID, value int) (Score, error) {
	criterion, err := rubric.Criterion(id)
	if err != nil {
		return Score{}, err
	}

	return NewScore(value, criterion.MaxScore())
}

func (s *Student) Evaluation(id CriterionID) (Evaluation, bool) {
	e, ok := s.evaluations[id]
	if !ok {
		return Evaluation{}, false
	}

	return *e, true
}
