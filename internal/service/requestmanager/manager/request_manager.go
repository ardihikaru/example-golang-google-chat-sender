// Package requestmanager handles all incoming printing request from the 3rd parties
package requestmanager

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ardihikaru/example-golang-google-chat-sender/internal/service/requestmanager/dto"
	"github.com/ardihikaru/example-golang-google-chat-sender/internal/storage/queue"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

// storage provides the interface for the functionality of MySQL Storage
type storage interface {
	Insert() error
}

// Service prepares the interfaces related with this minio service
type Service struct {
	log *logger.Logger

	// storage extends storage interface{} to implement storage.Storage
	// this variable should to be closed from modification: Open/Closed Principle (OCP)
	storage storage
}

// NewService creates a new user service
func NewService(log *logger.Logger, storage *queue.Store) *Service {
	return &Service{
		log:     log,
		storage: storage, // TODO: implements this
	}
}

// ExtractRmqMessage builds message
func (svc *Service) ExtractRmqMessage(rawMsg []byte, rawHeaders amqp.Table) (*dto.Message, error) {
	var err error

	// unmarshal body
	// FYI: we will always expect that the data type is a JSON
	var msg dto.Message
	if err = msg.ExtractMessage(rawMsg, rawHeaders); err != nil {
		svc.log.Error(err.Error())
		return nil, err
	}

	return &msg, nil
}

// Store stores to database
func (svc *Service) Store(message *dto.Message) error {
	// TODO: implements code here

	return svc.storage.Insert()
}
