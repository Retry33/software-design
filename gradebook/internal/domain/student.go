package domain

type Student struct {
	id   StudentID
	name string
}

func NewStudent(id StudentID, name string) (*Student, error) {
	normalizedName, ok := normalizeName(name)
	if !ok {
		return nil, ErrInvalidStudentName
	}

	return &Student{id: id, name: normalizedName}, nil
}

func (s *Student) Rename(name string) error {
	normalizedName, ok := normalizeName(name)
	if !ok {
		return ErrInvalidStudentName
	}

	s.name = normalizedName

	return nil
}

func (s *Student) ID() StudentID {
	return s.id
}

func (s *Student) Name() string {
	return s.name
}
