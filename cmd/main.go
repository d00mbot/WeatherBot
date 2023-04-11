package main

import (
	"subscription-bot/cmd/app"

	log "github.com/sirupsen/logrus"

	_ "time/tzdata"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Run(). Error: '%v'", err)
	}
}
