package rmqbuilder

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Builder defines the builder function
type Builder struct {
	Channel *amqp.Channel
}

// MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
}

// Message is the amqp request to publish
type Message struct {
	RoutingKey string
	Priority   uint8
	Body       MessageBody
}

// BuildStringBody builds string-typed body
func (b *Builder) BuildStringBody(body string) (*[]byte, error) {
	bytesBody := []byte(body)

	return &bytesBody, nil
}

// BuildRMQMessage builds a message for rabbitMQ
func (b *Builder) BuildRMQMessage(rmqRoute string, priority uint8, data interface{}) (*Message, error) {
	// Convert configData to string
	jsonBytesStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	//// test unmarshal
	//var dest map[string]interface{}
	//if err = json.Unmarshal(jsonBytesStr, &dest); err != nil {
	//	panic(err)
	//}
	//if err != nil {
	//	return nil, err
	//}

	// build rabbitMQ message
	msgBody := MessageBody{
		Data: jsonBytesStr,
	}
	message := &Message{
		RoutingKey: rmqRoute,
		Priority:   priority,
		Body:       msgBody,
	}

	return message, nil
}

// BuildMessageHeaders builds message headers for rabbitMQ
func (b *Builder) BuildMessageHeaders(rawHeaders map[string]interface{}) (*amqp.Table, error) {
	headers := amqp.Table{}

	for key, header := range rawHeaders {
		headers[key] = header
	}

	return &headers, nil
}
