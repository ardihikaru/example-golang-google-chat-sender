package dto

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/enum/rabbitmq"
)

func TestExtractMessageHeaders(t *testing.T) {
	var err error
	var message Message
	var want string
	var got string

	// builds message headers (amqp table)
	amqpHeaders := amqp.Table{}
	amqpHeaders[rmqenum.HeaderSvcId] = "service-id-123"
	amqpHeaders[rmqenum.HeaderRequeueStrategy] = "REQUEUE"
	amqpHeaders[rmqenum.HeaderMessageType] = rmqenum.Json
	amqpHeaders[rmqenum.HeaderTargetRouteKey] = "route-123"

	// builds message body
	rawMsgBytes := []byte("{\"extension\":\"pdf\",\"file_url\":\"ardi/sample.json\",\"printer_model\":\"Zebra\"}")

	err = message.ExtractMessage(rawMsgBytes, amqpHeaders)
	if err != nil {
		t.Errorf("test failed: build message failed")
	}

	// test service ID value
	want = amqpHeaders[rmqenum.HeaderSvcId].(string)
	got = message.Headers.SvcId
	if got != want {
		t.Errorf("Test fail! want: '%s', got: '%s'", want, got)
	}

	// test request strategy value
	want = amqpHeaders[rmqenum.HeaderRequeueStrategy].(string)
	got = message.Headers.RequeueStrategy
	if got != want {
		t.Errorf("Test fail! want: '%s', got: '%s'", want, got)
	}

	// test message type value
	want = amqpHeaders[rmqenum.HeaderMessageType].(string)
	got = message.Headers.MessageType
	if got != want {
		t.Errorf("Test fail! want: '%s', got: '%s'", want, got)
	}

	// test target route key value
	want = amqpHeaders[rmqenum.HeaderTargetRouteKey].(string)
	got = message.Headers.TargetRouteKey
	if got != want {
		t.Errorf("Test fail! want: '%s', got: '%s'", want, got)
	}
}
