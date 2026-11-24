package console

import (
	"errors"

	"github.com/Retry33/software-design/gradebook/internal/domain"
)

var errorMessages = []struct {
	target error
	text   string
}{
	{domain.ErrInvalidStudentID, "ID студента не может быть пустым"},
	{domain.ErrInvalidCriterionID, "ID критерия не может быть пустым"},
	{domain.ErrInvalidStudentName, "имя студента не может быть пустым"},
	{domain.ErrInvalidCriterionName, "название критерия не может быть пустым"},
	{domain.ErrInvalidMaxScore, "максимальный балл должен быть положительным"},
	{domain.ErrWeightOutOfRange, "вес должен быть от 1 до 100"},
	{domain.ErrWeightsExceedLimit, "сумма весов критериев превысит 100"},
	{domain.ErrScoreOutOfRange, "балл вне диапазона критерия"},
	{domain.ErrCriterionNotFound, "такого критерия нет в рубрикаторе"},
	{domain.ErrCriterionAlreadyExists, "критерий с таким ID уже есть"},
	{domain.ErrStudentNotFound, "студент не найден"},
	{domain.ErrStudentAlreadyExists, "студент с таким ID уже есть"},
	{domain.ErrEvaluationExists, "оценка по этому критерию уже выставлена, используйте «Исправить оценку»"},
	{domain.ErrEvaluationNotFound, "оценки по этому критерию нет"},
	{errNotANumber, "нужно ввести целое число"},
}

func describe(err error) string {
	for _, m := range errorMessages {
		if errors.Is(err, m.target) {
			return m.text
		}
	}

	return err.Error()
}
