package publisher

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/rabbitmq/builder"
)

type Publisher struct {
	log         *logger.Logger
	channel     *amqp.Channel
	queue       *amqp.Queue
	exchange    string
	connTimeout time.Duration
}

// NewPublisher creates a new publisher
func NewPublisher(log *logger.Logger, channel *amqp.Channel, exchange string, connTimeout time.Duration,
	queue *amqp.Queue) *Publisher {
	return &Publisher{
		log:         log,
		channel:     channel,
		queue:       queue,
		exchange:    exchange,
		connTimeout: connTimeout,
	}
}

// Publish publishes a request to the amqp queue
func (pub *Publisher) Publish(msg *rmqbuilder.Message, headers *amqp.Table, routingKey string) error {
	// enriches publication timeout
	ctx, cancel := context.WithTimeout(context.Background(), pub.connTimeout)
	defer cancel()

	err := pub.channel.PublishWithContext(ctx,
		//"", // exchange
		pub.exchange, // exchange
		//pub.queue.Name, // routing key
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Headers:     *headers,
			Body:        msg.Body.Data,
		})
	if err != nil {
		return err
	}

	pub.log.Debug("message has been published")

	return nil
}
