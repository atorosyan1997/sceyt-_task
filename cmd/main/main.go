package main

import (
	"sceyt_task/internal/app"
	"sceyt_task/internal/config"
)

// @title User-Server
// @version 1.0
// @description User Server for Sceyt test task.

// @host localhost:8080
// @BasePath /

// @in header
func main() {
	app.Run(config.ServerAddr, config.ServerPort)
}
