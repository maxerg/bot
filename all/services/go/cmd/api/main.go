package main

import (
	"log"

	"bot/internal/app"
	"bot/internal/config"
)

func main() {
	cfg := config.MustLoad()
	a := app.New(cfg)

	log.Println("BUILD_MARKER=AUTH_LOG_V1") 

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
