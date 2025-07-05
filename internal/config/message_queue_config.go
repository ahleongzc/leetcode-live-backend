package config

import (
	"strconv"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/util"
)

const (
	MESSAGE_QUEUE_DEV_HOST string = "amqp://guest:guest@localhost:5672/"
)

type MessageQueueConfig struct {
	Host                  string
	Queues                []string
	ReconnectionDelay     time.Duration
	ReinitializationDelay time.Duration
	ResendDelay           time.Duration
}

func LoadMessageQueueConfig() (*MessageQueueConfig, error) {
	host := util.GetEnvOr(common.MESSAGE_QUEUE_HOST_KEY, MESSAGE_QUEUE_DEV_HOST)

	reconnectionDelay := 5 * time.Second
	reinitializationDelay := 2 * time.Second
	resendDelay := 5 * time.Second

	reconnectionDelaySecValue, err := strconv.Atoi(util.GetEnvOr(common.MESSAGE_QUEUE_RECONNECT_DELAY_SEC_KEY, "5"))
	if nil == err {
		reconnectionDelay = time.Duration(reconnectionDelaySecValue) * time.Second
	}

	reinitializationDelaySecValue, err := strconv.Atoi(util.GetEnvOr(common.MESSAGE_QUEUE_REINIT_DELAY_SEC_KEY, "2"))
	if nil == err {
		reinitializationDelay = time.Duration(reinitializationDelaySecValue) * time.Second
	}

	resendDelaySecValue, err := strconv.Atoi(util.GetEnvOr(common.MESSAGE_QUEUE_RESEND_DELAY_SEC_KEY, "5"))
	if nil == err {
		resendDelay = time.Duration(resendDelaySecValue) * time.Second
	}

	// Define the queues here
	queues := []string{
		common.REVIEW_QUEUE,
	}

	return &MessageQueueConfig{
		Host:                  host,
		Queues:                queues,
		ReconnectionDelay:     reconnectionDelay,
		ReinitializationDelay: reinitializationDelay,
		ResendDelay:           resendDelay,
	}, nil
}
