package entity

import "github.com/ahleongzc/leetcode-live-backend/internal/util"

type Interview struct {
	Base
	UserID                uint
	QuestionID            uint
	StartTimestampMS      *int64
	ReviewID              *uint
	EndTimestampMS        *int64
	Token                 *string
	QuestionAttemptNumber uint
}

func (i *Interview) GetReviewID() uint {
	if i == nil || i.ReviewID == nil {
		return 0
	}

	return util.FromPtr(i.ReviewID)
}

func (i *Interview) GetStartTimesampS() int64 {
	if i == nil || i.StartTimestampMS == nil {
		return 0
	}
	return util.MillisToSeconds(util.FromPtr(i.StartTimestampMS))
}

func (i *Interview) GetEndTimestampS() int64 {
	if i == nil || i.EndTimestampMS == nil {
		return 0
	}
	return util.MillisToSeconds(util.FromPtr(i.EndTimestampMS))
}
