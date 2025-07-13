package entity

type Question struct {
	Base
	ExternalID  string `gorm:"index"`
	Description string
}

func NewQuestion() *Question {
	return &Question{}
}

func (q *Question) SetExternalID(externalID string) *Question {
	if q == nil {
		return nil
	}
	q.SetExternalID(externalID)
	return q
}

func (q *Question) SetDescription(description string) *Question {
	if q == nil {
		return nil
	}
	q.SetDescription(description)
	return q
}
