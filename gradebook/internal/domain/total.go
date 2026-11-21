package domain

type Total struct {
	value float64
}

func (t Total) Value() float64 {
	return t.value
}

type TotalCalculator struct {
	policy TotalPolicy
}

func NewTotalCalculator(policy TotalPolicy) *TotalCalculator {
	return &TotalCalculator{policy: policy}
}

func (c *TotalCalculator) Total(rubric *Rubric, student *Student) Total {
	criteria := rubric.Criteria()
	lines := make([]ScoredCriterion, 0, len(criteria))
	for _, criterion := range criteria {
		evaluation, graded := student.Evaluation(criterion.ID())
		lines = append(lines, ScoredCriterion{
			criterion: criterion,
			score:     evaluation.Score(),
			graded:    graded,
		})
	}

	return Total{value: c.policy.Total(lines)}
}
