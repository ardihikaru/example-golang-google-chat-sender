package taskmanagerutility

import (
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

// Utility is the interface
type Utility struct {
	Log *logger.Logger
}

// WindowsCmdUtility is the executable functions in within Windows command
type WindowsCmdUtility struct {
	Log *logger.Logger
}

// LinuxCmdUtility is the executable functions in within Linux command
type LinuxCmdUtility struct {
	Log *logger.Logger
}
