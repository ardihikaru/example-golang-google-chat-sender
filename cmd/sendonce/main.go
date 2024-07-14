// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions
// and limitations under the License.
package main

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	chat "google.golang.org/api/chat/v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardihikaru/example-golang-google-chat-sender/internal/application"
	"github.com/ardihikaru/example-golang-google-chat-sender/internal/config"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	e "github.com/ardihikaru/example-golang-google-chat-sender/pkg/utils/error"
)

func main() {
	var err error

	// loads config
	cfg, err := config.Get()
	if err != nil {
		e.FatalOnError(err, "failed to load config")
	}

	// validates configuration
	err = cfg.Validate()
	if err != nil {
		e.FatalOnError(err, "failed to load config")
	}

	// configures logger
	log, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		e.FatalOnError(err, "failed to prepare the logger")
	}

	err = os.Setenv("SERVICE_ACCOUNT_PATH", cfg.Google.ServiceAccount)

	// Setup client to write messages to chat.google.com
	client := getOauthClient(os.Getenv("SERVICE_ACCOUNT_PATH"))
	chatService, err := chat.New(client)
	if err != nil {
		e.FatalOnError(err, "failed to create chat service")
	}

	msgSvc := application.BuildSender(log, chatService)
	schedulerSvc := application.BuildScheduler(log, msgSvc)

	dateStr := "2024-07-14T14:51:00+07:00"
	for _, rawSpaceName := range cfg.Google.Spaces {
		spaceName := fmt.Sprintf("spaces/%s", rawSpaceName)
		log.Debug(fmt.Sprintf("space name: %s", spaceName))

		err = schedulerSvc.AddJobOneTimeTask(spaceName, dateStr)
		if err != nil {
			log.Warn(fmt.Sprintf("failed to send message to space [%s]: %s", spaceName, err.Error()))
		}
	}

	// starts the scheduler
	log.Debug("starting the scheduler")
	schedulerSvc.Start()

	dateStrAgain := "2024-07-14T14:52:00+07:00"
	for _, rawSpaceName := range cfg.Google.Spaces {
		spaceName := fmt.Sprintf("spaces/%s", rawSpaceName)
		log.Debug(fmt.Sprintf("space name: %s", spaceName))

		err = schedulerSvc.AddJobOneTimeTask(spaceName, dateStrAgain)
		if err != nil {
			log.Warn(fmt.Sprintf("failed to send message to space [%s]: %s", spaceName, err.Error()))
			e.FatalOnError(err, fmt.Sprintf("failed to send message to space [%s]", spaceName))
		}
	}

	go schedulerSvc.DoPeriodicLogging()

	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// shutdowns the application
	<-c
	log.Info("gracefully shutting down the scheduler")
	err = schedulerSvc.Shutdown()
	if err != nil {
		// handle error
		e.FatalOnError(err, "failed to showdown the scheduler app")
	}
	os.Exit(0)
}

// getOauthClient gets google OAuth client
func getOauthClient(serviceAccountKeyPath string) *http.Client {
	ctx := context.Background()
	data, err := os.ReadFile(serviceAccountKeyPath)
	if err != nil {
		e.FatalOnError(err, "failed to read service account file")
	}
	creds, err := google.CredentialsFromJSON(
		ctx,
		data,
		"https://www.googleapis.com/auth/chat.bot",
	)
	if err != nil {
		e.FatalOnError(err, "failed to fetch google creds")
	}
	return oauth2.NewClient(ctx, creds.TokenSource)
}
