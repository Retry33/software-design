package console

import "github.com/Retry33/software-design/gradebook/internal/application"

const demoExportFile = "demo-totals.csv"

func (c *Console) runDemo() error {
	students, err := c.app.Students()
	if err != nil {
		return err
	}
	criteria, err := c.app.Criteria()
	if err != nil {
		return err
	}
	if len(students) > 0 || len(criteria) > 0 {
		c.println("Демо запускается только на пустом журнале: перезапустите программу.")
		return nil
	}

	c.println("\n1) Заводим студентов и рубрикатор")
	for _, cmd := range []application.AddStudentCommand{
		{ID: "s1", Name: "Иван Петров"},
		{ID: "s2", Name: "Анна Иванова"},
		{ID: "s3", Name: "Олег Смирнов"},
	} {
		if _, err := c.app.AddStudent(cmd); err != nil {
			return err
		}
	}
	for _, cmd := range []application.AddCriterionCommand{
		{ID: "hw", Name: "Домашние задания", MaxScore: 10, Weight: 30},
		{ID: "test", Name: "Контрольная", MaxScore: 10, Weight: 30},
		{ID: "exam", Name: "Экзамен", MaxScore: 20, Weight: 40},
	} {
		if _, err := c.app.AddCriterion(cmd); err != nil {
			return err
		}
	}
	if err := c.listStudents(); err != nil {
		return err
	}
	if err := c.listCriteria(); err != nil {
		return err
	}

	c.println("\n2) Выставляем оценки")
	for _, cmd := range []application.GradeCommand{
		{StudentID: "s1", CriterionID: "hw", Score: 8},
		{StudentID: "s1", CriterionID: "test", Score: 9},
		{StudentID: "s1", CriterionID: "exam", Score: 15},
		{StudentID: "s2", CriterionID: "hw", Score: 10},
		{StudentID: "s2", CriterionID: "exam", Score: 18},
		{StudentID: "s3", CriterionID: "test", Score: 6},
	} {
		if err := c.app.Grade(cmd); err != nil {
			return err
		}
		c.printf("  %s: %s = %d\n", cmd.StudentID, cmd.CriterionID, cmd.Score)
	}

	c.println("\n3) Модель не дает нарушить правила")
	c.showRule("12 баллов из 10",
		c.app.Grade(application.GradeCommand{StudentID: "s3", CriterionID: "hw", Score: 12}))
	c.showRule("оценка по критерию вне рубрикатора",
		c.app.Grade(application.GradeCommand{StudentID: "s3", CriterionID: "lab", Score: 5}))
	c.showRule("вторая оценка по тому же критерию",
		c.app.Grade(application.GradeCommand{StudentID: "s1", CriterionID: "hw", Score: 7}))

	c.println("\n4) Итоги (считаются из текущих оценок и рубрикатора)")
	if err := c.listTotals(); err != nil {
		return err
	}

	c.println("\n5) Экспорт итогов")
	if err := c.app.ExportTotals(demoExportFile); err != nil {
		return err
	}
	c.printf("Итоги выгружены в %s.\n", demoExportFile)

	return nil
}

func (c *Console) showRule(rule string, err error) {
	if err == nil {
		c.printf("  ✗ %s: правило не сработало!\n", rule)
		return
	}
	c.printf("ok %s: отклонено (%s)\n", rule, describe(err))
}
