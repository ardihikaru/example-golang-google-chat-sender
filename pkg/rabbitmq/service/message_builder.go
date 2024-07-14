package rmqservice

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/rabbitmq/builder"
)

// BuildStringBody builds string body
func (svc *Service) BuildStringBody(body string) (*[]byte, error) {
	return svc.builder.BuildStringBody(body)
}

// BuildMessage builds rabbitMQ message
func (svc *Service) BuildMessage(rmqRoute string, priority uint8, data interface{}) (*rmqbuilder.Message, error) {
	return svc.builder.BuildRMQMessage(rmqRoute, priority, data)
}

// BuildMessageHeaders builds rabbitMQ message headers
func (svc *Service) BuildMessageHeaders(rawHeaders map[string]interface{}) (*amqp.Table, error) {
	return svc.builder.BuildMessageHeaders(rawHeaders)
}
