package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/Retry33/software-design/gradebook/internal/application"
)

const defaultExportFile = "totals.csv"

var (
	errInputClosed = errors.New("input closed")
	errNotANumber  = errors.New("not a number")
)

type Console struct {
	app *application.Gradebook
	in  *bufio.Scanner
	out io.Writer
}

func New(app *application.Gradebook, in io.Reader, out io.Writer) *Console {
	return &Console{app: app, in: bufio.NewScanner(in), out: out}
}

type menuItem struct {
	title string
	run   func() error
}

func (c *Console) menu() []menuItem {
	return []menuItem{
		{"Добавить студента", c.addStudent},
		{"Список студентов", c.listStudents},
		{"Добавить критерий", c.addCriterion},
		{"Рубрикатор", c.listCriteria},
		{"Выставить оценку", c.grade},
		{"Исправить оценку", c.correct},
		{"Отменить оценку", c.cancel},
		{"Оценки студента", c.listEvaluations},
		{"Итог студента", c.showTotal},
		{"Итоги всех студентов", c.listTotals},
		{"Экспорт итогов в CSV", c.exportTotals},
		{"Демо-сценарий", c.runDemo},
	}
}

func (c *Console) Run() error {
	menu := c.menu()
	for {
		c.printMenu(menu)
		choice, err := c.ask("Пункт меню")
		if errors.Is(err, errInputClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		if choice == "0" {
			c.println("До свидания!")
			return nil
		}

		n, convErr := strconv.Atoi(choice)
		if convErr != nil || n < 1 || n > len(menu) {
			c.println("Нет такого пункта меню.")
			continue
		}

		err = menu[n-1].run()
		if errors.Is(err, errInputClosed) {
			return nil
		}
		if err != nil {
			c.println("Ошибка: " + describe(err))
		}
	}
}

func (c *Console) printMenu(menu []menuItem) {
	c.println("")
	c.println("=== ПОС: преподавательская оценка студентов ===")
	for i, item := range menu {
		c.printf("%2d. %s\n", i+1, item.title)
	}
	c.println(" 0. Выход")
}

func (c *Console) addStudent() error {
	id, err := c.ask("ID студента")
	if err != nil {
		return err
	}
	name, err := c.ask("Имя")
	if err != nil {
		return err
	}

	view, err := c.app.AddStudent(application.AddStudentCommand{ID: id, Name: name})
	if err != nil {
		return err
	}
	c.printf("Добавлен студент %s (%s).\n", view.Name, view.ID)

	return nil
}

func (c *Console) listStudents() error {
	students, err := c.app.Students()
	if err != nil {
		return err
	}
	if len(students) == 0 {
		c.println("Студентов пока нет.")
		return nil
	}

	w := c.table("ID", "Имя")
	for _, s := range students {
		fmt.Fprintf(w, "%s\t%s\n", s.ID, s.Name)
	}

	return w.Flush()
}

func (c *Console) addCriterion() error {
	id, err := c.ask("ID критерия")
	if err != nil {
		return err
	}
	name, err := c.ask("Название")
	if err != nil {
		return err
	}
	maxScore, err := c.askInt("Максимальный балл")
	if err != nil {
		return err
	}
	weight, err := c.askInt("Вес, %")
	if err != nil {
		return err
	}

	view, err := c.app.AddCriterion(application.AddCriterionCommand{
		ID: id, Name: name, MaxScore: maxScore, Weight: weight,
	})
	if err != nil {
		return err
	}
	c.printf("Добавлен критерий %s (%s): максимум %d, вес %d%%.\n", view.Name, view.ID, view.MaxScore, view.Weight)

	return nil
}

func (c *Console) listCriteria() error {
	criteria, err := c.app.Criteria()
	if err != nil {
		return err
	}
	if len(criteria) == 0 {
		c.println("Рубрикатор пуст.")
		return nil
	}

	total := 0
	w := c.table("ID", "Название", "Максимум", "Вес, %")
	for _, cr := range criteria {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", cr.ID, cr.Name, cr.MaxScore, cr.Weight)
		total += cr.Weight
	}
	if err := w.Flush(); err != nil {
		return err
	}
	c.printf("Сумма весов: %d%%\n", total)

	return nil
}

func (c *Console) grade() error {
	cmd, err := c.askGradeCommand()
	if err != nil {
		return err
	}
	if err := c.app.Grade(cmd); err != nil {
		return err
	}
	c.printf("Студенту %s выставлено %d по критерию %s.\n", cmd.StudentID, cmd.Score, cmd.CriterionID)

	return nil
}

func (c *Console) correct() error {
	cmd, err := c.askGradeCommand()
	if err != nil {
		return err
	}
	if err := c.app.Correct(cmd); err != nil {
		return err
	}
	c.printf("Оценка студента %s по критерию %s исправлена на %d.\n", cmd.StudentID, cmd.CriterionID, cmd.Score)

	return nil
}

func (c *Console) cancel() error {
	studentID, err := c.ask("ID студента")
	if err != nil {
		return err
	}
	criterionID, err := c.ask("ID критерия")
	if err != nil {
		return err
	}

	cmd := application.CancelEvaluationCommand{StudentID: studentID, CriterionID: criterionID}
	if err := c.app.Cancel(cmd); err != nil {
		return err
	}
	c.printf("Оценка студента %s по критерию %s отменена.\n", studentID, criterionID)

	return nil
}

func (c *Console) listEvaluations() error {
	studentID, err := c.ask("ID студента")
	if err != nil {
		return err
	}
	evaluations, err := c.app.Evaluations(studentID)
	if err != nil {
		return err
	}
	if len(evaluations) == 0 {
		c.println("Оценок пока нет.")
		return nil
	}

	w := c.table("Критерий", "Название", "Балл")
	for _, e := range evaluations {
		fmt.Fprintf(w, "%s\t%s\t%d/%d\n", e.CriterionID, e.CriterionName, e.Score, e.MaxScore)
	}

	return w.Flush()
}

func (c *Console) showTotal() error {
	studentID, err := c.ask("ID студента")
	if err != nil {
		return err
	}
	total, err := c.app.Total(studentID)
	if err != nil {
		return err
	}
	c.printf("Итог %s (%s): %.2f\n", total.StudentName, total.StudentID, total.Total)

	return nil
}

func (c *Console) listTotals() error {
	totals, err := c.app.Totals()
	if err != nil {
		return err
	}
	if len(totals) == 0 {
		c.println("Студентов пока нет.")
		return nil
	}

	w := c.table("ID", "Имя", "Итог")
	for _, t := range totals {
		fmt.Fprintf(w, "%s\t%s\t%.2f\n", t.StudentID, t.StudentName, t.Total)
	}

	return w.Flush()
}

func (c *Console) exportTotals() error {
	path, err := c.ask("Файл (Enter — " + defaultExportFile + ")")
	if err != nil {
		return err
	}
	if path == "" {
		path = defaultExportFile
	}
	if err := c.app.ExportTotals(path); err != nil {
		return err
	}
	c.printf("Итоги выгружены в %s.\n", path)

	return nil
}

func (c *Console) ask(prompt string) (string, error) {
	c.printf("%s: ", prompt)
	if !c.in.Scan() {
		if err := c.in.Err(); err != nil {
			return "", err
		}
		return "", errInputClosed
	}

	return strings.TrimSpace(c.in.Text()), nil
}

func (c *Console) askInt(prompt string) (int, error) {
	raw, err := c.ask(prompt)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%q: %w", raw, errNotANumber)
	}

	return n, nil
}

func (c *Console) askGradeCommand() (application.GradeCommand, error) {
	studentID, err := c.ask("ID студента")
	if err != nil {
		return application.GradeCommand{}, err
	}
	criterionID, err := c.ask("ID критерия")
	if err != nil {
		return application.GradeCommand{}, err
	}
	score, err := c.askInt("Балл")
	if err != nil {
		return application.GradeCommand{}, err
	}

	return application.GradeCommand{StudentID: studentID, CriterionID: criterionID, Score: score}, nil
}

func (c *Console) table(headers ...string) *tabwriter.Writer {
	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	return w
}

func (c *Console) printf(format string, args ...any) {
	fmt.Fprintf(c.out, format, args...)
}

func (c *Console) println(s string) {
	fmt.Fprintln(c.out, s)
}
