package main

import (
	"linkshortener/internal/config"
	"linkshortener/server"
)

func main() {
	config.LoadConfig(".")
	server.Init()
}
