package main

import (
	"subscription-bot/app"
	_ "time/tzdata"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Run(). Error: '%v'", err)
	}
}
