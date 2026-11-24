package domain

import "strings"

type CriterionID struct {
	value string
}

func NewCriterionID(value string) (CriterionID, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return CriterionID{}, ErrInvalidCriterionID
	}
	return CriterionID{value: trimmedValue}, nil
}

func (id CriterionID) Value() string {
	return id.value
}
