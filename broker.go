package broker

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// eventQueue describes an event queue system to manage
// the queue messages.
type eventQueue struct {
	conn   *amqp.Connection
	events map[string]*handler
}

// event defines an entity to hold all information
// about the event.
type handler struct {
	HandleFn handlerFunc
	Options  dialOptions
}

// handlerFunc is an alias about the function who
// describes the main handle to manage the queue messages.
type handlerFunc = func(req *Request) (*Response, error)

// Reques holds the queue message.
type Request struct {
	Body io.Reader
}

type Response struct{}

// DialOption defines a function to set up the event
// queue system resources, such as queue name or route keys.
type DialOption func(*dialOptions)

// dialOptions represents the options config about
// the rabbitmq resources to declare its entities.
type dialOptions struct {
	QueueName string
	RouteKey  string
}

type MessageSystem interface {
	OpenChannel() (EventQueueChannel, error)
	PublishEvent(event, routeKey string, msg io.Reader, ch EventQueueChannel) error
	Close() error
}

// NewEventQueue creates new instance of eventQueue
// and open new rabbitmq connection.
func NewEventQueue(host, port, user, password string) (*eventQueue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", user, password, host, port))
	if err != nil {
		return nil, err
	}
	return &eventQueue{
		conn:   conn,
		events: map[string]*handler{},
	}, nil
}

type EventQueueChannel interface {
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error
	Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
	Close() error
}

func (eq *eventQueue) OpenChannel() (EventQueueChannel, error) {
	ch, err := eq.conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// PublishEvent publishes an event to exchange entity using an routeKey and msg value.
func (eq *eventQueue) PublishEvent(event, routeKey string, msg io.Reader, ch EventQueueChannel) error {
	if err := ch.ExchangeDeclare(
		event,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	body, err := ioutil.ReadAll(msg)
	if err != nil {
		return err
	}
	err = ch.Publish(
		event,
		routeKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to Publish")
	}
	return nil
}

func (eq *eventQueue) Close() error {
	return eq.conn.Close()
}
