package model

type UserProfile struct {
	Username                string `json:"username"`
	Email                   string `json:"email"`
	RemainingInterviewCount uint   `json:"remaining_interview_count"`
	InterviewDurationS      uint   `json:"interview_duration_s"`
}

func NewUserProfile() *UserProfile {
	return &UserProfile{}
}

func (u *UserProfile) SetUsername(username string) *UserProfile {
	if u == nil {
		return nil
	}
	u.Username = username
	return u
}

func (u *UserProfile) SetEmail(email string) *UserProfile {
	if u == nil {
		return nil
	}
	u.Email = email
	return u
}

func (u *UserProfile) SetRemainingInterviewCount(count uint) *UserProfile {
	if u == nil {
		return nil
	}
	u.RemainingInterviewCount = count
	return u
}

func (u *UserProfile) SetInterviewDurationS(durationSeconds uint) *UserProfile {
	if u == nil {
		return nil
	}
	u.InterviewDurationS = durationSeconds
	return u
}
