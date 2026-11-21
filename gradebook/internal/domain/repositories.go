package domain

type RubricRepository interface {
	Get() (*Rubric, error)
	Save(r *Rubric) error
}

type StudentRepository interface {
	ByID(id StudentID) (*Student, error)
	All() ([]*Student, error)
	Save(s *Student) error
}
