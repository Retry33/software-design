package domain

const (
	minWeight = 1
	maxWeight = 100
)

type Weight struct {
	value int
}

func NewWeight(value int) (Weight, error) {
	if value < minWeight || value > maxWeight {
		return Weight{}, ErrWeightOutOfRange
	}
	return Weight{value: value}, nil
}

func (w Weight) Value() int {
	return w.value
}
