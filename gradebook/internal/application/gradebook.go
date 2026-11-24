package application

type Gradebook struct {
	students *StudentService
	rubric   *RubricService
	grading  *GradingService
	reports  *ReportService
}

func NewGradebook(
	students *StudentService,
	rubric *RubricService,
	grading *GradingService,
	reports *ReportService,
) *Gradebook {
	return &Gradebook{students: students, rubric: rubric, grading: grading, reports: reports}
}

func (g *Gradebook) AddStudent(cmd AddStudentCommand) (StudentView, error) {
	return g.students.AddStudent(cmd)
}

func (g *Gradebook) Students() ([]StudentView, error) {
	return g.students.Students()
}

func (g *Gradebook) AddCriterion(cmd AddCriterionCommand) (CriterionView, error) {
	return g.rubric.AddCriterion(cmd)
}

func (g *Gradebook) Criteria() ([]CriterionView, error) {
	return g.rubric.Criteria()
}

func (g *Gradebook) Grade(cmd GradeCommand) error {
	return g.grading.Grade(cmd)
}

func (g *Gradebook) Correct(cmd GradeCommand) error {
	return g.grading.Correct(cmd)
}

func (g *Gradebook) Cancel(cmd CancelEvaluationCommand) error {
	return g.grading.Cancel(cmd)
}

func (g *Gradebook) Evaluations(studentID string) ([]EvaluationView, error) {
	return g.grading.Evaluations(studentID)
}

func (g *Gradebook) Total(studentID string) (TotalView, error) {
	return g.reports.Total(studentID)
}

func (g *Gradebook) Totals() ([]TotalView, error) {
	return g.reports.Totals()
}

func (g *Gradebook) ExportTotals(destination string) error {
	return g.reports.ExportTotals(destination)
}
