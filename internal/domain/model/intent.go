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

func (i *IntentDetail) GetIntentWithHighestConfidence() Intent {
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

func (i *IntentDetail) GetIntentWithHighestConfidenceWithScoreOutOf100() (Intent, float64) {
	intentWithHighestScore := i.GetIntentWithHighestConfidence()
	if score, ok := i.Mapping[intentWithHighestScore]; ok {
		return intentWithHighestScore, score * 100
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

func (i *IntentDetail) Exists() bool {
	return i != nil
}
