package domain

import "math"

type ScoredCriterion struct {
	criterion *Criterion
	score     Score
	graded    bool
}

func (sc ScoredCriterion) Graded() bool { return sc.graded }

func (sc ScoredCriterion) Weight() float64 { return float64(sc.criterion.Weight().Value()) }

func (sc ScoredCriterion) Contribution() float64 {
	if !sc.graded {
		return 0
	}

	return float64(sc.score.Value()) / float64(sc.criterion.MaxScore().Value()) * sc.Weight()
}

type TotalPolicy interface {
	Total(lines []ScoredCriterion) float64
}

type StrictPolicy struct{}

func (StrictPolicy) Total(lines []ScoredCriterion) float64 {
	total := 0.0
	for _, l := range lines {
		total += l.Contribution()
	}

	return total
}

type GradedOnlyPolicy struct{}

func (GradedOnlyPolicy) Total(lines []ScoredCriterion) float64 {
	contribution, weight := 0.0, 0.0
	for _, l := range lines {
		if !l.Graded() {
			continue
		}
		contribution += l.Contribution()
		weight += l.Weight()
	}
	if weight == 0 {
		return 0
	}

	return contribution / weight * maxTotalWeight
}

type RoundedPolicy struct {
	inner  TotalPolicy
	places int
}

func NewRoundedPolicy(inner TotalPolicy, places int) RoundedPolicy {
	return RoundedPolicy{inner: inner, places: max(places, 0)}
}

func (p RoundedPolicy) Total(lines []ScoredCriterion) float64 {
	factor := math.Pow(10, float64(p.places))

	return math.Round(p.inner.Total(lines)*factor) / factor
}
