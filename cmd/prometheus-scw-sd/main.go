package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/promhippie/prometheus-scw-sd/pkg/command"
)

func main() {
	if env := os.Getenv("PROMETHEUS_SCW_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	if err := command.Run(); err != nil {
		os.Exit(1)
	}
}
