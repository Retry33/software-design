package application

type TotalsExporter interface {
	Export(destination string, totals []TotalView) error
}
