package domain

import "errors"

var (
	ErrScoreOutOfRange        = errors.New("score out of range")
	ErrInvalidMaxScore        = errors.New("max score must be positive")
	ErrWeightOutOfRange       = errors.New("weight out of range")
	ErrInvalidStudentID       = errors.New("invalid StudentID")
	ErrInvalidCriterionID     = errors.New("invalid CriterionID")
	ErrInvalidCriterionName   = errors.New("criterion name must not be empty")
	ErrInvalidStudentName     = errors.New("student name must not be empty")
	ErrCriterionNotFound      = errors.New("criterion not found in rubric")
	ErrCriterionAlreadyExists = errors.New("criterion already exists in rubric")
	ErrWeightsExceedLimit     = errors.New("total weight of criteria exceeds 100")
	ErrEvaluationExists       = errors.New("evaluation for this criterion already exists")
	ErrEvaluationNotFound     = errors.New("evaluation for this criterion not found")
	ErrStudentNotFound        = errors.New("student not found")
	ErrStudentAlreadyExists   = errors.New("student already exists")
)
