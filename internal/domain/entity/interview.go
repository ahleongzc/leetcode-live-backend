package entity

import (
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

type Interview struct {
	Base
	UserID               uint
	QuestionID           uint
	Code                 string
	StartTimestampMS     *int64
	ReviewID             *uint
	EndTimestampMS       *int64
	Token                *string
	QuestionAttemptCount uint
	SetupCount           uint
	Ongoing              bool
	Abandoned            bool
	AbandonedTimestampMS *int64
	ElapsedTimeS         uint
	AllocatedDurationS   uint
}

func NewInterview() *Interview {
	return &Interview{}
}

// The chances of this happening is low
// When the client starts an interview, it is a two step process
// 1. Setup interview (HTTP)
// 2. Join interview (Websocket)
// The client should join immediately after the initial set up, because the join interview will reset the setup count
func (i *Interview) ExceedSetupCountThreshold() bool {
	return i.SetupCount >= 3
}

func (i *Interview) GetTimeRemainingS() uint {
	return i.AllocatedDurationS - i.ElapsedTimeS - uint(util.MillisToSeconds(time.Now().UnixMilli()-i.UpdateTimestampMS))
}

func (i *Interview) TimesUp() bool {
	if i == nil {
		return true
	}
	timeRemainingS := i.GetTimeRemainingS()
	return timeRemainingS <= 0
}

func (i *Interview) SetAllocatedDurationS(durationSeconds uint) *Interview {
	if i == nil {
		return nil
	}
	i.AllocatedDurationS = durationSeconds
	return i
}

func (i *Interview) IncrementSetupCount() *Interview {
	if i == nil {
		return nil
	}
	i.SetupCount++
	return i
}

func (i *Interview) SetUserID(userID uint) *Interview {
	if i == nil {
		return nil
	}
	i.UserID = userID
	return i
}

func (i *Interview) SetQuestionID(questionID uint) *Interview {
	if i == nil {
		return nil
	}
	i.QuestionID = questionID
	return i
}

func (i *Interview) SetToken(token string) *Interview {
	if i == nil {
		return nil
	}
	i.Token = util.ToPtr(token)
	return i
}

func (i *Interview) SetQuestionAttemptCount(count uint) *Interview {
	if i == nil {
		return nil
	}
	i.QuestionAttemptCount = count
	return i
}

func (i *Interview) SetSetupCount(count uint) *Interview {
	if i == nil {
		return nil
	}
	i.SetupCount = count
	return i
}

func (i *Interview) SetReviewID(reviewID uint) *Interview {
	if i == nil {
		return nil
	}
	i.ReviewID = util.ToPtr(reviewID)
	return i
}

func (i *Interview) ResetSetupCount() {
	if i == nil {
		return
	}
	i.SetSetupCount(0)
}

func (i *Interview) Pause() *Interview {
	if i == nil {
		return nil
	}
	i.Ongoing = false
	i.UpdateElapsedTimeS()
	return i
}

func (i *Interview) SetOngoing() *Interview {
	if i == nil {
		return nil
	}
	i.Ongoing = true
	return i
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

// Note that the logic for abandon and update seems similar,
// but the key difference is that abandoning an interview doesn't update the elapsed duration
func (i *Interview) Abandon() *Interview {
	if i == nil {
		return nil
	}

	i.ConsumeToken()

	i.Ongoing = false
	i.EndTimestampMS = util.ToPtr(time.Now().UnixMilli())

	i.AbandonedTimestampMS = util.ToPtr(time.Now().UnixMilli())
	i.Abandoned = true

	return i
}

func (i *Interview) TokenExists() bool {
	if i == nil {
		return false
	}

	return i.Token != nil
}

func (i *Interview) ConsumeToken() *Interview {
	if i == nil {
		return nil
	}
	i.Token = nil
	return i
}

func (i *Interview) End() *Interview {
	if i == nil {
		return nil
	}

	i.Ongoing = false
	i.EndTimestampMS = util.ToPtr(time.Now().UnixMilli())

	i.UpdateElapsedTimeS()

	return i
}

func (i *Interview) Start() *Interview {
	if i == nil {
		return nil
	}
	i.StartTimestampMS = util.ToPtr(time.Now().UnixMilli())
	return i
}

func (i *Interview) UpdateElapsedTimeS() {
	if i == nil {
		return
	}

	durationS := util.MillisToSeconds(time.Now().UnixMilli() - i.UpdateTimestampMS)
	i.ElapsedTimeS += uint(durationS)
}

func (i *Interview) HasEnded() bool {
	if i == nil {
		return false
	}
	return i.EndTimestampMS != nil && !i.Ongoing
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

func (i *Interview) Exists() bool {
	return i != nil
}

func (i *Interview) GetToken() string {
	if i == nil || i.Token == nil {
		return ""
	}
	return util.FromPtr(i.Token)
}

func (i *Interview) ReviewExists() bool {
	if i == nil {
		return false
	}
	return i.ReviewID != nil
}
