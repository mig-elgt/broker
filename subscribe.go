package broker

import (
	"bytes"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

func (eq *eventQueue) subscribe(event string, handle *handler, done chan<- doneErrCh) {
	// Open new channel
	ch, err := eq.conn.Channel()
	if err != nil {
		done <- doneErrCh{event, err}
		return
	}
	defer ch.Close()

	// Exchange Declaration
	if err := ch.ExchangeDeclare(
		event,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		done <- doneErrCh{event, err}
		return
	}

	// Queue Dleclaration
	_, err = ch.QueueDeclare(
		handle.Options.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		done <- doneErrCh{event, err}
		return
	}

	// Binding Declaration
	if err := ch.QueueBind(handle.Options.QueueName, handle.Options.RouteKey, event, false, nil); err != nil {
		done <- doneErrCh{event, err}
		return
	}

	msgs, err := ch.Consume(
		handle.Options.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		done <- doneErrCh{event, err}
		return
	}
	logrus.Infof("waiting to receive messages from %v event", event)
	for d := range msgs {
		_, err := handle.HandleFn(&Request{Body: bytes.NewBuffer(d.Body)})
		if err != nil {
			logrus.WithField("event", event).Errorf("failed to perform message: %v; waiting 5 seconds to try again...", err)
			time.Sleep(10 * time.Second)
			d.Nack(false, true)
			continue
		}
		// TODO: validate response to publish another event
		d.Ack(false)
	}
	done <- doneErrCh{event, errors.New("channel is closed")}
}
