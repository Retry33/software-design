package domain

const maxTotalWeight = 100

type Rubric struct {
	criteria []*Criterion
}

func NewRubric() *Rubric {
	return &Rubric{}
}

func (r *Rubric) AddCriterion(c *Criterion) error {
	if c == nil || c.ID().Value() == "" {
		return ErrInvalidCriterionID
	}
	if _, ok := r.find(c.ID()); ok {
		return ErrCriterionAlreadyExists
	}
	if r.totalWeight()+c.Weight().Value() > maxTotalWeight {
		return ErrWeightsExceedLimit
	}

	r.criteria = append(r.criteria, c)

	return nil
}

func (r *Rubric) Criterion(id CriterionID) (*Criterion, error) {
	c, ok := r.find(id)
	if !ok {
		return nil, ErrCriterionNotFound
	}

	return c, nil
}

func (r *Rubric) Criteria() []*Criterion {
	return append([]*Criterion(nil), r.criteria...)
}

func (r *Rubric) find(id CriterionID) (*Criterion, bool) {
	for _, c := range r.criteria {
		if c.ID() == id {
			return c, true
		}
	}

	return nil, false
}

func (r *Rubric) totalWeight() int {
	total := 0
	for _, c := range r.criteria {
		total += c.Weight().Value()
	}

	return total
}
