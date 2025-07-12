package model

import "fmt"

type Intent string

const (
	CANDIDATE_EXPLANATION Intent = "explanation"
	OTHERS                Intent = "others"
	DEFAULT               Intent = "default"
)

type IntentDetail struct {
	Mapping map[Intent]float64
}

func NewIntentDetail() *IntentDetail {
	return &IntentDetail{
		Mapping: make(map[Intent]float64),
	}
}

func (i *IntentDetail) GetIntentWithHighestConfidenceScore() Intent {
	var highestScore float64
	var intentWithHighestScore Intent

	for k, v := range i.Mapping {
		if v > float64(highestScore) {
			intentWithHighestScore = k
			highestScore = v
		}
	}

	return intentWithHighestScore
}

func (i *IntentDetail) GetIntentWithHighestConfidenceScoreWithScore() (Intent, float64) {
	intentWithHighestScore := i.GetIntentWithHighestConfidenceScore()
	if score, ok := i.Mapping[intentWithHighestScore]; ok {
		return intentWithHighestScore, score
	}

	return DEFAULT, 0
}

func (i *IntentDetail) String() string {
	res := ""
	for k, v := range i.Mapping {
		res += fmt.Sprintf("%s has confidence score of %f\n", k, v)
	}

	return res
}
