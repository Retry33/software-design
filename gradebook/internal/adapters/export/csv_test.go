package export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Retry33/software-design/gradebook/internal/application"
)

var sampleTotals = []application.TotalView{
	{StudentID: "s1", StudentName: "Иван Петров", Total: 76.5},
	{StudentID: "s2", StudentName: "Doe, John", Total: 30},
}

const sampleCSV = "student_id,student_name,total\n" +
	"s1,Иван Петров,76.50\n" +
	"s2,\"Doe, John\",30.00\n"

func TestWriteTotals(t *testing.T) {
	// Arrange
	var buf bytes.Buffer

	// Act
	err := writeTotals(&buf, sampleTotals)

	// Assert
	if err != nil {
		t.Fatalf("writeTotals: %v", err)
	}
	if buf.String() != sampleCSV {
		t.Errorf("csv =\n%s\nwant\n%s", buf.String(), sampleCSV)
	}
}

func TestCSVExporterWritesFile(t *testing.T) {
	// Arrange
	path := filepath.Join(t.TempDir(), "totals.csv")

	// Act
	err := NewCSVExporter().Export(path, sampleTotals)

	// Assert
	if err != nil {
		t.Fatalf("Export: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != sampleCSV {
		t.Errorf("file =\n%s\nwant\n%s", data, sampleCSV)
	}
}
