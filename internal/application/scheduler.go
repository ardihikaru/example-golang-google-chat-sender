package application

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	e "github.com/ardihikaru/example-golang-google-chat-sender/pkg/utils/error"
)

type Scheduler struct {
	gocron.Scheduler
	senderSvc *Sender
	log       *logger.Logger
	messages  []string
}

type data struct {
	hour   uint
	minute uint
	second uint
}

// BuildScheduler builds scheduler
func BuildScheduler(log *logger.Logger, msgSvc *Sender) *Scheduler {
	// create a scheduler
	cronScheduler, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		e.FatalOnError(err, "failed to build a new scheduler")
	}

	return &Scheduler{
		Scheduler: cronScheduler,
		senderSvc: msgSvc,
		log:       log,
		messages:  GetScheduleList(),
	}
}

// Start starts the scheduler
func (sch *Scheduler) Start() {
	// start the scheduler
	sch.Scheduler.Start()
}

// extractScheduleData extracts schedule data
func (sch *Scheduler) extractScheduleData(schedule string) (*data, error) {
	var err error
	var dt data

	rawData := strings.Split(schedule, ":")

	hour, err := strconv.ParseUint(rawData[0], 10, 32)
	if err != nil {
		return nil, err
	}
	dt.hour = uint(hour)

	minute, err := strconv.ParseUint(rawData[1], 10, 32)
	if err != nil {
		return nil, err
	}
	dt.minute = uint(minute)

	second, err := strconv.ParseUint(rawData[2], 10, 64)
	if err != nil {
		return nil, err
	}
	dt.second = uint(second)

	return &dt, nil
}

// getRandMessage gets a random message
func (sch *Scheduler) getRandMessage() string {
	//msg := sch.messages[len(sch.messages)-1]
	//sch.log.Debug(fmt.Sprintf("selected message: %s", msg))
	//
	//// remove selected message
	//sch.messages = sch.messages[:len(sch.messages)-1]
	randKey := randRange(1, 5)

	return GetScheduleMsg()[randKey]
}

// randRange gets a random value
func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// BuildAtTime builds gocron AtTime object
func (sch *Scheduler) BuildAtTime(schedule string) gocron.AtTime {
	dt, err := sch.extractScheduleData(schedule)
	if err != nil {
		e.FatalOnError(err, "build schedule data failed")
	}
	sch.log.Debug(fmt.Sprintf("adding a cronjob message at %d:%d:%d", dt.hour, dt.minute, dt.second))

	// builds gocron AtTime
	return gocron.NewAtTime(dt.hour, dt.minute, dt.second)
}

// AddJob adds a new job
func (sch *Scheduler) AddJob(spaceName string, atTime gocron.AtTime) error {
	start := time.Now()
	//start = start.Add(1 * time.Second) // starts it 1 second after the scheduler has been started
	start = start.Add(24 * time.Hour) // starts it 1 second after the scheduler has been started

	// if messages emptied already, reject any job request
	if len(sch.messages) == 0 {
		return fmt.Errorf("no more message to generate")
	}

	// selects the message to be added to the job
	message := sch.getRandMessage()

	job, err := sch.NewJob(
		gocron.WeeklyJob(
			1,
			//gocron.NewWeekdays(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday),
			gocron.NewWeekdays(time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday),
			gocron.NewAtTimes(
				atTime,
				atTime,
				//gocron.NewAtTime(12, 0, 30),
				//gocron.NewAtTime(23, 0, 0),
			),
		),
		gocron.NewTask(
			func(spaceName, message string) {
				// executes the action here
				sch.senderSvc.sendMessage(spaceName, message)
			},
			spaceName, message,
		),
		gocron.WithStartAt(
			gocron.WithStartDateTime(start),
		),
	)
	if err != nil {
		// handle error
		sch.log.Warn("failed to add a job to the scheduler")
		return err
	}

	// each job has a unique id
	sch.log.Info(fmt.Sprintf("added job ID: %s", job.ID()))

	return nil
}

// AddJobOneTimeTask adds a new job to run a one time task
func (sch *Scheduler) AddJobOneTimeTask(spaceName, dateStr string) error {
	// if messages emptied already, reject any job request
	if len(sch.messages) == 0 {
		return fmt.Errorf("no more message to generate")
	}

	// selects the message to be added to the job
	message := sch.getRandMessage()

	runAt, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return fmt.Errorf("failed to decode time")
	}

	job, err := sch.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(runAt)),
		gocron.NewTask(
			func(spaceName, message string) {
				// executes the action here
				sch.senderSvc.sendMessage(spaceName, message)
			},
			spaceName, message,
		),
	)
	if err != nil {
		// handle error
		sch.log.Warn("failed to add a job to the scheduler")
		return err
	}

	// each job has a unique id
	sch.log.Info(fmt.Sprintf("added job ID: %s", job.ID()))

	return nil
}

// DoPeriodicLogging executes log message periodically
func (sch *Scheduler) DoPeriodicLogging() {
	for {
		sch.log.Debug("I am logging periodically ... in a background process")

		sch.log.Debug("list of jobs:")
		for _, job := range sch.Jobs() {
			sch.log.Debug(fmt.Sprintf("job ID: %v", job.ID()))
		}

		time.Sleep(2 * time.Second)
	}
}
