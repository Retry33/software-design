package domain

type Criterion struct {
	id       CriterionID
	name     string
	maxScore MaxScore
	weight   Weight
}

func NewCriterion(id CriterionID, name string, maxScore MaxScore, w Weight) (*Criterion, error) {
	if id.Value() == "" {
		return nil, ErrInvalidCriterionID
	}
	if maxScore.Value() <= 0 {
		return nil, ErrInvalidMaxScore
	}
	if w.Value() < minWeight || w.Value() > maxWeight {
		return nil, ErrWeightOutOfRange
	}
	normalizedName, ok := normalizeName(name)
	if !ok {
		return nil, ErrInvalidCriterionName
	}
	return &Criterion{
		id:       id,
		name:     normalizedName,
		maxScore: maxScore,
		weight:   w,
	}, nil
}

func (c *Criterion) Rename(name string) error {
	normalizedName, ok := normalizeName(name)
	if !ok {
		return ErrInvalidCriterionName
	}
	c.name = normalizedName
	return nil
}

func (c *Criterion) ID() CriterionID {
	return c.id
}

func (c *Criterion) Name() string {
	return c.name
}

func (c *Criterion) MaxScore() MaxScore {
	return c.maxScore
}

func (c *Criterion) Weight() Weight {
	return c.weight
}
