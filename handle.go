package broker

import "fmt"

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
