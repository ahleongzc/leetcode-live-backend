//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ahleongzc/leetcode-live-backend/cmd/app"

	"github.com/google/wire"
)

func InitializeApplication() (*app.Application, error) {
	wire.Build(
		// Application
		app.NewApplication,
	)
	return &app.Application{}, nil
}
