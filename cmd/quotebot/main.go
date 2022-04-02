package main

import (
	"github.com/Depal/quotebot/internal/entry"
	"log"
)

func main() {
	app := entry.Initialize()

	err := app.Setup()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = app.Teardown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	app.Start()
}
