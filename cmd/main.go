package main

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/wire"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app, err := wire.InitializeApplication()
	if err != nil {
		panic(err)
	}

	app.StartBackgroundTasks()

	serverConfig := config.LoadServerConfig()

	server := &http.Server{
		Addr:         serverConfig.Address,
		Handler:      app.Handler(),
		IdleTimeout:  serverConfig.IdleTimeout,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
