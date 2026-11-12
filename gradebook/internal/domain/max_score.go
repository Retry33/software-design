package domain

type MaxScore struct {
	value int
}

func NewMaxScore(value int) (MaxScore, error) {
	if value <= 0 {
		return MaxScore{}, ErrInvalidMaxScore
	}
	return MaxScore{value: value}, nil
}

func (ms MaxScore) Value() int {
	return ms.value
}
