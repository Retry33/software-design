package domain

type Evaluation struct {
	criterionID CriterionID
	score       Score
}

func newEvaluation(criterionID CriterionID, score Score) *Evaluation {
	return &Evaluation{criterionID: criterionID, score: score}
}

func (e Evaluation) CriterionID() CriterionID {
	return e.criterionID
}

func (e Evaluation) Score() Score {
	return e.score
}
