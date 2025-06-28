package model

type InterviewMessageType string

const (
	URL InterviewMessageType = "url"
)

type InterviewMessage struct {
	Type    InterviewMessageType
	Content string
}
