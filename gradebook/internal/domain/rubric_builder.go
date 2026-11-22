package domain

import "fmt"

type RubricBuilder struct {
	rubric *Rubric
	err    error
}

func NewRubricBuilder() *RubricBuilder {
	return &RubricBuilder{rubric: NewRubric()}
}

func (b *RubricBuilder) Criterion(id, name string, maxScore, weight int) *RubricBuilder {
	if b.err != nil {
		return b
	}
	if err := b.add(id, name, maxScore, weight); err != nil {
		b.err = fmt.Errorf("criterion %q: %w", id, err)
	}

	return b
}

func (b *RubricBuilder) Build() (*Rubric, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.rubric, nil
}

func (b *RubricBuilder) add(id, name string, maxScore, weight int) error {
	criterionID, err := NewCriterionID(id)
	if err != nil {
		return err
	}
	ms, err := NewMaxScore(maxScore)
	if err != nil {
		return err
	}
	w, err := NewWeight(weight)
	if err != nil {
		return err
	}
	criterion, err := NewCriterion(criterionID, name, ms, w)
	if err != nil {
		return err
	}

	return b.rubric.AddCriterion(criterion)
}

func NewRubricBuilderFrom(rubric *Rubric) *RubricBuilder {
	return &RubricBuilder{rubric: &Rubric{criteria: rubric.Criteria()}}
}
