package main

import (
	"bytes"
	"log"
	"time"

	"github.com/mig-elgt/broker"
)

func main() {
	eq, err := broker.NewEventQueue("localhost", "5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	defer eq.Close()

	ch, err := eq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	for {
		if err := eq.PublishEvent("tracking", "new_locations", bytes.NewBufferString(`{"user": 100}`), ch); err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
}
