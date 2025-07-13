package entity

type Review struct {
	Base
	Score    uint
	Passed   bool
	Feedback string
}

func NewReview() *Review {
	return &Review{}
}

func (r *Review) SetScore(score uint) *Review {
	if r == nil {
		return nil
	}
	r.Score = score
	return r
}

func (r *Review) SetPassed(passed bool) *Review {
	if r == nil {
		return nil
	}
	r.Passed = passed
	return r
}

func (r *Review) SetFeedback(feedback string) *Review {
	if r == nil {
		return nil
	}
	r.Feedback = feedback
	return r
}

func (r *Review) Exists() bool {
	return r != nil
}
