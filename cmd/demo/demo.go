package main

import (
	"log"

	"github.com/joncalhoun/multireaderdemos"
)

func main() {
	err := multireaderdemos.JSONLogsWithoutMultiReader()
	if err != nil {
		log.Fatal(err)
	}
	err = multireaderdemos.JSONLogsWithMultiReader()
	if err != nil {
		log.Fatal(err)
	}
}
