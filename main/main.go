package main

import (
	"log"
)

func main() {
	app := application{}
	app.LoadAndRoute()

	log.Fatal(app.server.ListenAndServe())
}
