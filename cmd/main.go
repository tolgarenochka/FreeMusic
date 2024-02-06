package main

import (
	"flag"

	"FreeMusic/internal/app"
)

var configPath *string

// @title FreeMusic
// @version 1.0
// @description API Server for FreeMusic Application

func init() {
	configPath = flag.String("config-path", "./configs/local_config.json", "path to config")
	flag.Parse()
}

func main() {
	app.Run(*configPath)
}
