package domain

import "strings"

type StudentID struct {
	value string
}

func NewStudentID(value string) (StudentID, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return StudentID{}, ErrInvalidStudentID
	}
	return StudentID{value: trimmedValue}, nil
}

func (id StudentID) Value() string {
	return id.value
}
