package model

type Delivery struct {
	Body []byte
	Acknowledger
}

type Acknowledger interface {
	Ack() error
	Nack(requeue bool) error
	Reject(requeue bool) error
}
