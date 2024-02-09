package scheduler

import (
	"fmt"
	"gameAppProject/param"
	"gameAppProject/service/matchingservice"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Config struct {
	MatchWaitedUsersIntervalInSeconds int `koanf:"match_waited_users_interval_in_seconds"`
}

type Scheduler struct {
	sch      *gocron.Scheduler
	matchSvc matchingservice.Service
	config   Config
}

func New(config Config, matchSvc matchingservice.Service) Scheduler {
	return Scheduler{
		config:   config,
		matchSvc: matchSvc,
		sch:      gocron.NewScheduler(time.UTC)}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	// TODO - add lock or distributed lock

	defer wg.Done()

	s.sch.Every(s.config.MatchWaitedUsersIntervalInSeconds).Second().Do(s.MatchWaitedUsers)

	s.sch.StartAsync()

	<-done
	// wait to finish job
	fmt.Println("stop scheduler..")
	s.sch.Stop()
}

func (s Scheduler) MatchWaitedUsers() {
	s.matchSvc.MatchWaitedUsers(param.MatchWaitedUsersRequest{})
}
