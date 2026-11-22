package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Retry33/software-design/gradebook/internal/adapters/console"
	"github.com/Retry33/software-design/gradebook/internal/adapters/export"
	"github.com/Retry33/software-design/gradebook/internal/adapters/memory"
	"github.com/Retry33/software-design/gradebook/internal/application"
	"github.com/Retry33/software-design/gradebook/internal/domain"
	"go.uber.org/dig"
)

func main() {
	policyName := flag.String("policy", "strict", "правило расчета итога: strict или graded-only")
	flag.Parse()
	if err := run(*policyName); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run(policyName string) error {
	var policy domain.TotalPolicy
	switch policyName {
	case "strict":
		policy = domain.StrictPolicy{}
	case "graded-only":
		policy = domain.GradedOnlyPolicy{}
	default:
		return fmt.Errorf("неизвестное правило %q, допустимы strict и graded-only", policyName)
	}

	container := dig.New()
	providers := []struct {
		constructor any
		options     []dig.ProvideOption
	}{
		{memory.NewStudentRepository, []dig.ProvideOption{dig.As(new(domain.StudentRepository))}},
		{memory.NewRubricRepository, []dig.ProvideOption{dig.As(new(domain.RubricRepository))}},
		{export.NewCSVExporter, []dig.ProvideOption{dig.As(new(application.TotalsExporter))}},
		{func() domain.TotalPolicy { return domain.NewRoundedPolicy(policy, 2) }, nil},
		{domain.NewTotalCalculator, nil},
		{application.NewStudentService, nil},
		{application.NewRubricService, nil},
		{application.NewGradingService, nil},
		{application.NewReportService, nil},
		{application.NewGradebook, nil},
		{func(app *application.Gradebook) *console.Console {
			return console.New(app, os.Stdin, os.Stdout)
		}, nil},
	}
	for _, provider := range providers {
		if err := container.Provide(provider.constructor, provider.options...); err != nil {
			return fmt.Errorf("регистрация зависимости: %w", err)
		}
	}
	return container.Invoke((*console.Console).Run)
}
