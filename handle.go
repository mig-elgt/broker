package broker

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// HandleEvent registers a new event with its handler. It set up the default
// options config about the queue and route key.
func (eq *eventQueue) HandleEvent(eventName string, handleFn handlerFunc, opts ...DialOption) {
	dopts := dialOptions{
		QueueName: fmt.Sprintf("%s.queue", eventName),
		RouteKey:  "",
	}
	for _, opt := range opts {
		opt(&dopts)
	}
	eq.events[eventName] = &handler{
		HandleFn: handleFn,
		Options:  dopts,
	}
}

func (eq *eventQueue) RunSubscriber(eventName string, handleFn handlerFunc, opts ...DialOption) {
	dopts := dialOptions{
		QueueName: fmt.Sprintf("%s.queue", eventName),
		RouteKey:  "",
	}
	for _, opt := range opts {
		opt(&dopts)
	}
	errorCh := make(chan doneErrCh)
	eventHandler := &handler{
		HandleFn: handleFn,
		Options:  dopts,
	}
	go eq.subscribe(eventName, eventHandler, errorCh)
	for {
		select {
		case e := <-errorCh:
			logrus.Errorf("subscriber to event %v died: %v", e.event, e.err)
			logrus.Infof("waiting 10 seconds to re-run subscriber to %v event", e.event)
			time.Sleep(10 * time.Second)
			// Re-run subscriber
			go eq.subscribe(eventName, eventHandler, errorCh)
		}
	}
}

// WithQueue sets new queue name.
func WithQueue(name string) DialOption {
	return func(do *dialOptions) {
		do.QueueName = name
	}
}

// WithRouteKey sets new route key name.
func WithRouteKey(name string) DialOption {
	return func(do *dialOptions) {
		do.RouteKey = name
	}
}
