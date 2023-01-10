Broker wraps RabbitMQ client to connect a set of queues using concurrency.

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
		if err := eq.PublishEvent("foobar", "route_key", bytes.NewBufferString("hello world"), ch); err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
}

```

# Subscriber Example

```go
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
	eq.HandleEvent("event_name", func(req *broker.Request) (*broker.Response, error) {
		type message struct {
			ID string `json:"id"`
		}
		m := message{}
		if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
			return nil, err
		}
		return &broker.Response{}, nil
	}, broker.WithQueue("queue_name"), broker.WithRouteKey("route_key"))

	<-eq.RunSubscribers()
	log.Fatalf("rabbitmq server connection is closed %v", err)
}
```

