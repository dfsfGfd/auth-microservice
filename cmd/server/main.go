package main

import (
	"fmt"
	"os"

	"auth-microservice/internal/di"
	"auth-microservice/pkg/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	app, err := di.InitializeApplication()
	if err != nil {
		return err
	}
	defer app.CleanUp()

	srv := server.NewServer(app.Config, app.Logger, app.RateLimiter, app.AuthHandler)
	return srv.Run()
}
