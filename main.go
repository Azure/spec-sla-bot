package main

import (
	"log"

	"github.com/Azure/spec-sla-bot/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
