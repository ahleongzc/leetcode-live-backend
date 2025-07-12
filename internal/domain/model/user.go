package model

type UserProfile struct {
	Username                string `json:"username"`
	Email                   string `json:"email"`
	RemainingInterviewCount uint   `json:"remaining_interview_count"`
	InterviewDurationS      uint   `json:"interview_duration_s"`
}

type UserProfileBuilder struct {
	username                string
	email                   string
	remainingInterviewCount uint
	interviewDurationS      uint
}

func NewUserProfileBuilder() *UserProfileBuilder {
	return &UserProfileBuilder{}
}

func (u *UserProfileBuilder) SetUsername(username string) *UserProfileBuilder {
	u.username = username
	return u
}

func (u *UserProfileBuilder) SetEmail(email string) *UserProfileBuilder {
	u.email = email
	return u
}

func (u *UserProfileBuilder) SetRemainingInterviewCount(count uint) *UserProfileBuilder {
	u.remainingInterviewCount = count
	return u
}

func (u *UserProfileBuilder) SetInterviewDurationS(duration uint) *UserProfileBuilder {
	u.interviewDurationS = duration
	return u
}

func (u *UserProfileBuilder) Build() *UserProfile {
	return &UserProfile{
		Username:                u.username,
		Email:                   u.email,
		RemainingInterviewCount: u.remainingInterviewCount,
		InterviewDurationS:      u.interviewDurationS,
	}
}
