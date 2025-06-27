package infra

import (
	"io"
	"os"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/util"
	"github.com/rs/zerolog"
)

func NewZerologLogger() zerolog.Logger {
	var writer io.Writer

	if util.IsDevEnv() {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
			NoColor:    false,
		}
	} else {
		writer = os.Stdout
	}

	return zerolog.New(writer).With().Timestamp().Logger()
}
