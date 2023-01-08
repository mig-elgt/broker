package broker

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type doneErrCh struct {
	event string
	err   error
}

func (eq *eventQueue) RunSubscribers() <-chan *amqp.Error {
	ch := make(chan *amqp.Error)
	go func() {
		// Create errors channel list
		errorsCh := []chan doneErrCh{}
		for i := 0; i < len(eq.events); i++ {
			errorsCh = append(errorsCh, make(chan doneErrCh))
		}
		// Run subscribers
		go func() {
			count := 0
			for event, handle := range eq.events {
				go eq.subscribe(event, handle, errorsCh[count])
				count++
			}
		}()
		// Listen errors
		for _, ch := range errorsCh {
			ch := ch
			go func() {
				for {
					select {
					case e := <-ch:
						logrus.Errorf("subscriber to event %v died: %v", e.event, e.err)
						logrus.Infof("waiting 10 seconds to re-run subscriber to %v event", e.event)
						time.Sleep(10 * time.Second)
						// Re-run subscriber
						go eq.subscribe(e.event, eq.events[e.event], ch)
					}
				}
			}()
		}
	}()
	return eq.conn.NotifyClose(ch)
}
