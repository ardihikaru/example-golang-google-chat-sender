package dto

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Message defines the whatsapp message
type Message struct {
	Body             MessageBody    `json:"body"`
	Headers          MessageHeaders `json:"headers"`
	SourceRoutingKey string         `json:"source_routing_key"`
}

// MessageBody defines th message body
type MessageBody struct {
	FileExt      string `json:"extension"`
	FilePath     string `json:"file_path"`
	PrinterModel string `json:"printer_model"`
}

// MessageHeaders defines th message headers
type MessageHeaders struct {
	SvcId           string `json:"svc_id"`
	RequeueStrategy string `json:"requeue_strategy"`
	MessageType     string `json:"message_type"`
}

// Validate validates message
func (msg *Message) Validate() error {
	// TODO: implements this code
	// 		 e.g. extension is not supported

	return nil
}

// ExtractMessage extracts message body
func (msg *Message) ExtractMessage(rawMsg []byte, rawHeaders amqp.Table) error {
	var err error

	// unmarshal body
	// FYI: we will always expect that the data type is a JSON
	if err = msg.extractMessageBody(rawMsg); err != nil {
		return fmt.Errorf("failed to build message body: %s", err.Error())
	}

	if err = msg.extractMessageHeaders(rawHeaders); err != nil {
		return fmt.Errorf("failed to build message header: %s", err.Error())
	}

	return nil
}

// extractMessageBody extracts message body
func (msg *Message) extractMessageBody(rawMsg []byte) error {
	var err error

	if err = json.Unmarshal(rawMsg, &msg.Body); err != nil {
		return err
	}

	return nil
}

// extractMessageHeaders extracts message headers
func (msg *Message) extractMessageHeaders(rawHeaders amqp.Table) error {
	// casts to bytes first
	headerBytes, _ := json.Marshal(rawHeaders)

	// to message headers
	err := json.Unmarshal(headerBytes, &msg.Headers)
	if err != nil {
		return err
	}

	return nil
}
