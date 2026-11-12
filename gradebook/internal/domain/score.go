package domain

type Score struct {
	value int
}

func NewScore(value int, maxScore MaxScore) (Score, error) {
	if maxScore.value <= 0 {
		return Score{}, ErrInvalidMaxScore
	}

	if value < 0 || value > maxScore.value {
		return Score{}, ErrScoreOutOfRange
	}
	return Score{value: value}, nil
}

func (s Score) Value() int {
	return s.value
}
