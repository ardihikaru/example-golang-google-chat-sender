// Package rmqservice provides all functionality and modules related with rabbitMQ
package rmqservice

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/config"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/rabbitmq/builder"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/rabbitmq/publisher"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/rabbitmq/subscriber"
	e "github.com/ardihikaru/example-golang-google-chat-sender/pkg/utils/error"
)

const (
	exchangeType = "topic"
)

type builder interface {
	BuildStringBody(body string) (*[]byte, error)
	BuildRMQMessage(rmqRoute string, priority uint8, data interface{}) (*rmqbuilder.Message, error)
	BuildMessageHeaders(rawHeaders map[string]interface{}) (*amqp.Table, error)
}

// Service is the connection created
type Service struct {
	Conn     *amqp.Connection
	Channel  *amqp.Channel
	Pub      *publisher.Publisher
	Sub      *subscriber.Subscriber
	log      *logger.Logger
	rmqCfg   *config.RabbitMQ
	exchange string
	err      chan error

	// builder extends builder interface{} to implement rmqbuilder.Builder
	// this variable should to be closed from modification: Open/Closed Principle (OCP)
	builder builder
}

var (
	connectionPool = make(map[string]*Service)
)

// NewService returns the new connection object
func NewService(rmqCfg config.RabbitMQ, market, buildMode string, log *logger.Logger) *Service {
	var err error

	// builds exchange
	exchange := buildExchange(rmqCfg.Exchange, market, buildMode, rmqCfg.WithBuildMode)

	// shows the built exchange
	log.Info("",
		zap.Bool("with build mode", rmqCfg.WithBuildMode),
		zap.String("exchange", exchange),
	)

	// assigns to the connection pool
	// when it tries to open the same pool again, returns the currently open connection instead of creating another one
	if svc, ok := connectionPool[rmqCfg.PoolName]; ok {
		return svc
	}

	// builds service
	svc := &Service{
		rmqCfg:   &rmqCfg,
		exchange: exchange,
		log:      log,
		builder:  nil,
		err:      make(chan error),
	}

	// initiates channel to the builder
	// at this point the value will be nil
	svc.builder = &rmqbuilder.Builder{
		Channel: svc.Channel,
	}

	// opens rabbitMQ connections
	if err = svc.Connect(); err != nil {
		e.FatalOnError(err, "failed to initialize rabbitMQ connection")
	}

	// stores to the pool
	connectionPool[rmqCfg.PoolName] = svc

	return svc
}

// buildExchange builds the exchange based on the config
// if withBuildMode=true => system will build the routing key with "<Market>_<BuildMode>_<RmqExchange>"
// if EnableBuildMode=false => system will build the routing key with "<Market>_<RmqExchange>"
func buildExchange(exchangePostfix, market, buildMode string, withBuildMode bool) string {
	if withBuildMode {
		return fmt.Sprintf("%s_%s_%s", market, buildMode, exchangePostfix)
	} else {
		return fmt.Sprintf("%s_%s", market, exchangePostfix)
	}
}

// Connect builds the rabbitMQ connection
func (svc *Service) Connect() error {
	var err error

	// build config opt
	cfgOpts := amqp.Config{
		Heartbeat: svc.rmqCfg.Heartbeat,
		Dial:      amqp.DefaultDial(svc.rmqCfg.ConnTimeout),
	}

	svc.Conn, err = amqp.DialConfig(svc.rmqCfg.URI, cfgOpts)
	if err != nil {
		return fmt.Errorf("error in creating rabbitmq connection with %s : %s", svc.rmqCfg.URI, err.Error())
	}
	go func() {
		// From the AMQP library docs:
		//   NotifyClose method on the amqp.Channel object registers a listener
		//   for when the server sends a channel or connection exception in the form
		//   of a Connection.Close or Channel.Close method
		<-svc.Conn.NotifyClose(make(chan *amqp.Error)) // listens to NotifyClose

		// at this point, instead of restart, we will stop the system and restart
		// until the problem with the rabbitMQ connection has been stabilized
		svc.log.Fatal("rabbitMQ connection has been closed")

		e.FatalOnError(err, fmt.Sprintf("rabbitMQ connection has been unexpectedly closed"))
	}()
	svc.Channel, err = svc.Conn.Channel()
	if err != nil {
		return fmt.Errorf("got error when opening the channel: %s", err)
	}

	// declares the exchange with topic type
	// TODO: sets the ACK into a manual mode
	// TODO: implements priority mode with a range between 1-10
	if err := svc.Channel.ExchangeDeclare(
		svc.exchange, // name
		exchangeType, // type
		false,        // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("error in declaring a new exchange: %s", err)
	}

	// binds route
	err = svc.Channel.QueueBind(
		svc.rmqCfg.QueueName,
		svc.rmqCfg.Route,
		svc.exchange,
		false,
		nil,
	)
	if err != nil {
		e.FatalOnError(err, "failed to bind a routing key")
	}

	svc.log.Info(fmt.Sprintf("listens to queue: %s", svc.rmqCfg.QueueName))
	svc.log.Info(fmt.Sprintf("binds route key: %s", svc.rmqCfg.Route))

	// builds and declares queue
	queue, err := svc.BuildQueue()
	if err != nil {
		e.FatalOnError(err, "failed to declare a queue")
	}

	// once the connection ready.
	// 1) sets the publisher (if enabled)
	svc.buildPublisher(&queue)

	// 2) once the connection ready, sets the subscriber (if enabled)
	svc.buildSubscriber(&queue)

	return nil
}

// buildSubscriber builds subscriber instance
func (svc *Service) buildSubscriber(queue *amqp.Queue) {
	if svc.rmqCfg.SubscriberEnabled {
		svc.Sub = subscriber.NewSubscriber(svc.log, svc.Channel, queue)
	}
}

// buildPublisher builds publisher instance
func (svc *Service) buildPublisher(queue *amqp.Queue) {
	if svc.rmqCfg.PublisherEnabled {
		svc.Pub = publisher.NewPublisher(
			svc.log,
			svc.Channel,
			svc.exchange,
			svc.rmqCfg.ConnTimeout,
			queue,
		)
	}
}

// BuildQueue declares message queue
func (svc *Service) BuildQueue() (amqp.Queue, error) {
	return svc.Channel.QueueDeclare(
		svc.rmqCfg.QueueName, // name
		svc.rmqCfg.Durable,   // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
}
