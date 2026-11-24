package application

import "github.com/Retry33/software-design/gradebook/internal/domain"

type ReportService struct {
	students   domain.StudentRepository
	rubrics    domain.RubricRepository
	calculator *domain.TotalCalculator
	exporter   TotalsExporter
}

func NewReportService(
	students domain.StudentRepository,
	rubrics domain.RubricRepository,
	calculator *domain.TotalCalculator,
	exporter TotalsExporter,
) *ReportService {
	return &ReportService{students: students, rubrics: rubrics, calculator: calculator, exporter: exporter}
}

func (s *ReportService) Total(studentID string) (TotalView, error) {
	student, err := loadStudent(s.students, studentID)
	if err != nil {
		return TotalView{}, err
	}
	rubric, err := s.rubrics.Get()
	if err != nil {
		return TotalView{}, err
	}

	return s.totalView(rubric, student), nil
}

func (s *ReportService) Totals() ([]TotalView, error) {
	students, err := s.students.All()
	if err != nil {
		return nil, err
	}
	rubric, err := s.rubrics.Get()
	if err != nil {
		return nil, err
	}

	totals := make([]TotalView, 0, len(students))
	for _, student := range students {
		totals = append(totals, s.totalView(rubric, student))
	}

	return totals, nil
}

func (s *ReportService) ExportTotals(destination string) error {
	totals, err := s.Totals()
	if err != nil {
		return err
	}

	return s.exporter.Export(destination, totals)
}

func (s *ReportService) totalView(rubric *domain.Rubric, student *domain.Student) TotalView {
	return TotalView{
		StudentID:   student.ID().Value(),
		StudentName: student.Name(),
		Total:       s.calculator.Total(rubric, student).Value(),
	}
}
