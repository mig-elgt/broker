package main

import (
	"encoding/json"
	"log"

	"github.com/mig-elgt/broker"
)

func main() {
	eq, err := broker.NewEventQueue("localhost", "5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	defer eq.Close()
	eq.HandleEvent("email.confirmation", func(req *broker.Request) (*broker.Response, error) {
		type message struct {
			Email string `json:"email"`
		}
		m := message{}
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			return nil, err
		}
		return &broker.Response{}, nil
	}, broker.WithQueue("emails_to_confirm"), broker.WithRouteKey("emails"))

	<-eq.RunSubscribers()
	log.Fatalf("rabbitmq server connection is closed %v", err)
}
