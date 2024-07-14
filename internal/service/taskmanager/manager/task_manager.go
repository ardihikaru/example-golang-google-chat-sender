// Package taskmanager handles captured tasks given by the request manager
package taskmanager

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ardihikaru/example-golang-google-chat-sender/internal/enum/operatingsystem"
	"github.com/ardihikaru/example-golang-google-chat-sender/internal/enum/printermodel"
	"github.com/ardihikaru/example-golang-google-chat-sender/internal/service/taskmanager/dto"
	"github.com/ardihikaru/example-golang-google-chat-sender/internal/service/taskmanager/utility"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

// utility provides the interface for the functionality of MinIO
type utility interface {
	PrintWithZebra(params *dto.ZebraParams) error
	PrintWithPos(params *dto.PosParams) error
	PrintWithDotMatrix(params *dto.DotMatrixParams) error
}

// Service prepares the interfaces related with this minio service
type Service struct {
	log *logger.Logger

	// utility extends utility interface{} to implement taskmanagerutility.Utility
	// this variable should to be closed from modification: Open/Closed Principle (OCP)
	utility utility
}

// NewService creates a new user service
func NewService(log *logger.Logger, osModel string) *Service {
	switch model := osModel; model {
	case operatingsystem.Windows:
		return buildServiceBasedOnWindows(log)
	case operatingsystem.Linux:
		return buildServiceBasedOnLinux(log)
	default:
		// sets default as Linux
		return buildServiceBasedOnLinux(log)
	}
}

func buildServiceBasedOnWindows(log *logger.Logger) *Service {
	log.Info("sets the printer command as Windows Command")
	return &Service{
		log: log,
		utility: &taskmanagerutility.WindowsCmdUtility{
			Log: log,
		},
	}
}

func buildServiceBasedOnLinux(log *logger.Logger) *Service {
	log.Info("sets default printer command as Linux Command")
	return &Service{
		log: log,
		utility: &taskmanagerutility.LinuxCmdUtility{
			Log: log,
		},
	}
}

// ExtractRmqMessage builds message
func (svc *Service) ExtractRmqMessage(rawMsg []byte, rawHeaders amqp.Table, routingKey string) (*dto.Message, error) {
	var err error

	// unmarshal body
	// FYI: we will always expect that the data type is a JSON
	var msg dto.Message
	if err = msg.ExtractMessage(rawMsg, rawHeaders); err != nil {
		svc.log.Error(err.Error())
		return nil, err
	}

	// captures source routing key
	msg.SourceRoutingKey = routingKey

	return &msg, nil
}

// Print prints message
func (svc *Service) Print(printerModel string, params interface{}) error {
	// TODO: implements more printer model here
	switch printerModel {
	case printermodel.Zebra:
		return svc.utility.PrintWithZebra(params.(*dto.ZebraParams))
	case printermodel.Pos:
		return svc.utility.PrintWithPos(params.(*dto.PosParams))
	case printermodel.DotMatrix:
		return svc.utility.PrintWithDotMatrix(params.(*dto.DotMatrixParams))
	default:
		// unknown printer model, returns an error message
		return fmt.Errorf("unknown printer model (%s)", printerModel)
	}
}
