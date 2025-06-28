package entity

type Role string

const (
	SYSTEM    Role = "system"
	USER      Role = "user"
	ASSISTANT Role = "assistant"
)

type Transcript struct {
	ID                 int
	Role               Role
	Content            string
	InterviewID        int
	CreatedTimestampMS int64
}
