package entity

import (
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type Interview struct {
	Base
	UserID                uint
	QuestionID            uint
	Code                  string
	StartTimestampMS      *int64
	ReviewID              *uint
	EndTimestampMS        *int64
	Token                 *string
	QuestionAttemptNumber uint
	SetUpCount            uint
	Ongoing               bool
	Abandoned             bool
	AbandonedTimestampMS  *int64
}

func (i *Interview) Pause() {
	if i == nil {
		return
	}
	i.Ongoing = false
}

func (i *Interview) SetOngoing() {
	if i == nil {
		return
	}
	i.Ongoing = true
}

func (i *Interview) IsUnstarted() bool {
	if i == nil {
		return false
	}

	return i.StartTimestampMS == nil
}

func (i *Interview) IsUnfinished() bool {
	if i == nil {
		return false
	}

	return i.StartTimestampMS != nil && i.EndTimestampMS == nil
}

func (i *Interview) Abandon() {
	if i == nil {
		return
	}
	i.AbandonedTimestampMS = util.ToPtr(time.Now().UnixMilli())
	i.Abandoned = true
	i.Ongoing = false
}

func (i *Interview) ConsumeToken() {
	if i == nil {
		return
	}
	i.Token = nil
}

func (i *Interview) End() {
	if i == nil {
		return
	}
	i.Ongoing = false
	i.EndTimestampMS = util.ToPtr(time.Now().UnixMilli())
}

func (i *Interview) Start() {
	if i == nil {
		return
	}
	i.StartTimestampMS = util.ToPtr(time.Now().UnixMilli())
}

func (i *Interview) HasStarted() bool {
	if i == nil {
		return false
	}
	return i.StartTimestampMS != nil
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
