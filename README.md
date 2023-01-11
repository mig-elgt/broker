Broker wraps Go AMQP 0.9.1 client (https://github.com/streadway/amqp) to expose an API in order to implement Event Driven Design pattern for Microservices Architecture.

# Install

```bash
go get github.com/mig-elgt/broker
```

# Publisher Example

```go
package main

import (
	"bytes"
	"log"
	"time"

	"github.com/mig-elgt/broker"
)

func main() {
	// Create Event Queue object
	eq, err := broker.NewEventQueue("localhost", "5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	defer eq.Close()
	// Open a new Channel
	ch, err := eq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	// Publish Event
	if err := eq.PublishEvent("foobar", "route_key", bytes.NewBufferString("hello world"), ch); err != nil { 
		log.Fatal(err)
	}
}

```

# Subscriber Example
Register your event handlers and create your function handler to perform a queue message.

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/mig-elgt/broker"
)

func main() {
	// Create Event Queue instance
	eq, err := broker.NewEventQueue("localhost", "5672", "guest", "guest")
	if err != nil {
		panic(err)
	}
	defer eq.Close()

	// Register Handle Events
	eq.HandleEvent("foo", func(req *broker.Request) (*broker.Response, error) {
	    // Add your code here
		return &broker.Response{}, nil
	}, broker.WithQueue("foo_queue"), broker.WithRouteKey("route_key"))

	eq.HandleEvent("bar", func(req *broker.Request) (*broker.Response, error) {
	    // Add your code here
		return &broker.Response{}, nil
	}, broker.WithQueue("bar_queue"), broker.WithRouteKey("route_key"))

	// Exec runners for each event
	<-eq.RunSubscribers()
	log.Fatalf("rabbitmq server connection is closed %v", err)
}

```
