package main

import (
	"account-management-service/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
