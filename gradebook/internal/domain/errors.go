package domain

import "errors"

var (
	ErrScoreOutOfRange  = errors.New("score out of range")
	ErrInvalidMaxScore  = errors.New("max score must be positive")
	ErrWeightOutOfRange = errors.New("weight out of range")
)
