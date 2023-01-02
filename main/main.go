package main

import (
	"log"
)

func main() {
	app := Application{}
	app.LoadAndRoute()

	log.Fatal(app.server.ListenAndServe())
}
