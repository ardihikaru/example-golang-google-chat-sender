package subscriber

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

// Subscriber defines the subscriber object
type Subscriber struct {
	Log     *logger.Logger
	Channel *amqp.Channel
	Queue   *amqp.Queue
}

// NewSubscriber creates a new publisher
func NewSubscriber(log *logger.Logger, channel *amqp.Channel, queue *amqp.Queue) *Subscriber {
	return &Subscriber{
		Log:     log,
		Channel: channel,
		Queue:   queue,
	}
}

// BuildMessageQueue builds consumable messages
func (sub *Subscriber) BuildMessageQueue() (<-chan amqp.Delivery, error) {
	if sub.Queue == nil {
		return nil, fmt.Errorf("message queue is not defined yet")
	}

	return sub.Channel.Consume(
		sub.Queue.Name, // queue
		"",             // consumer
		false,          // auto-ack; in our case, always set to false
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
}
