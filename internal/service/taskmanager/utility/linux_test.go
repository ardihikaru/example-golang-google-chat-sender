package taskmanagerutility

import (
	"testing"

	"github.com/ardihikaru/example-golang-google-chat-sender/internal/service/taskmanager/dto"
)

func TestPrintWithZebra(t *testing.T) {
	var err error

	//log := logger.NewTest(t)
	////logger.Debug("This is a debug message, debug messages will be logged as well")
	////logger.Info("Logs will not be shown during normal test execution")
	////logger.Warn("You can see all log messages of a successful run by passing the -v flag")
	////logger.Error("Additionally the entire log output for a specific unit test will be visible when a test fails")
	//
	//osModel := "linux"
	//log.Info(fmt.Sprintf("os model: %s", osModel))
	params := buildParams()
	//log.Info(fmt.Sprintf("printer IP: %s", params.IpPrinter))

	utility := LinuxCmdUtility{}

	err = utility.PrintWithZebra(params)
	if err != nil {
		t.Errorf("Test fail! failed to execute the print() method: %s", err.Error())
	}
}

// buildParams builds zebra printer parameters
func buildParams() *dto.ZebraParams {
	return &dto.ZebraParams{
		Device:    "",
		IpPrinter: "",
		Data:      nil,
	}
}
