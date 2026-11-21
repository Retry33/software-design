package application

type AddStudentCommand struct {
	ID   string
	Name string
}

type AddCriterionCommand struct {
	ID       string
	Name     string
	MaxScore int
	Weight   int
}

type GradeCommand struct {
	StudentID   string
	CriterionID string
	Score       int
}

type CancelEvaluationCommand struct {
	StudentID   string
	CriterionID string
}
