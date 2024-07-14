package application

import (
	"fmt"

	"google.golang.org/api/chat/v1"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	e "github.com/ardihikaru/example-golang-google-chat-sender/pkg/utils/error"
)

type Sender struct {
	*chat.SpacesMessagesService
	log *logger.Logger
}

// BuildSender builds sender
func BuildSender(log *logger.Logger, chatService *chat.Service) *Sender {
	return &Sender{
		SpacesMessagesService: chat.NewSpacesMessagesService(chatService),
		log:                   log,
	}
}

// sendMessage sends a message to the designated space
//
//	FYI: it will fail it this chat app has not been added to the designated space yet!
func (s *Sender) sendMessage(spaceName, rawMsg string) {
	var err error

	// builds message
	msg := &chat.Message{
		Text: fmt.Sprintf(rawMsg),
	}

	_, err = s.Create(spaceName, msg).Do()
	if err != nil {
		e.FatalOnError(err, "failed to create message")
	}
}
