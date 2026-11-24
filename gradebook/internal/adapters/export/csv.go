// Package export writes reports to files (secondary adapter).
package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Retry33/software-design/gradebook/internal/application"
)

var _ application.TotalsExporter = (*CSVExporter)(nil)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter { return &CSVExporter{} }

func (e *CSVExporter) Export(destination string, totals []application.TotalView) (err error) {
	f, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("create %s: %w", destination, err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()

	return writeTotals(f, totals)
}

func writeTotals(w io.Writer, totals []application.TotalView) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"student_id", "student_name", "total"}); err != nil {
		return err
	}
	for _, t := range totals {
		row := []string{t.StudentID, t.StudentName, strconv.FormatFloat(t.Total, 'f', 2, 64)}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()

	return cw.Error()
}
