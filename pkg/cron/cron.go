package cron

import "github.com/robfig/cron/v3"

type CronJob struct {
	Name     string
	Schedule string
	Func     func()
}

type CronService struct {
	cron *cron.Cron
}

var cronService *CronService

func NewCronService() *CronService {
	if cronService == nil {
		cronService = &CronService{
			cron: cron.New(),
		}
	}
	return cronService
}

func (cs *CronService) RegisterJob(job CronJob) error {
	_, err := cs.cron.AddFunc(job.Schedule, job.Func)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CronService) Start() {
	cs.cron.Start()
}

func (cs *CronService) Stop() {
	cs.cron.Stop()
}
