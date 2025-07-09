package model

type Intent string

const (
	NO_INTENT             Intent = "nil"
	HINT_REQUEST          Intent = "hint"
	CLARIFICATION_REQUEST Intent = "clarification"
	END_REQUEST           Intent = "end"
)
