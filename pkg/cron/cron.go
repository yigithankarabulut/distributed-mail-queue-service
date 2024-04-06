package cron

import "time"

//s, err := gocron.NewScheduler()
//if err != nil {
//apiserv.logger.Error("error creating scheduler", "error", err)
//}
//
//cronJob := &CronJob{
//JobName:   "FindUnprocessedTasksAndEnqueue",
//Scheduler: s,
//Task:      gocron.NewTask(taskService.FindUnprocessedTasksAndEnqueue, context.TODO()),
//Duration:  gocron.DurationJob(10 * time.Second),
//}
//jobID := cronJob.AddJob()
//if jobID == uuid.Nil {
//apiserv.logger.Error("error adding job")
//}
//log.Infof("job id: %s", jobID)
//s.Start()

type Scheduler interface {
	NewJob(duration Duration, task Task) (Job, error)
	Start()
}

type Job interface {
	ID() string
}

type Task func()

type Duration struct {
	time.Duration
}

func DurationJob(d time.Duration) Duration {
	return Duration{Duration: d}
}

type CronJob struct {
	JobName   string
	Scheduler Scheduler
	Task      Task
	Duration  Duration
}

//func (c *CronJob) AddJob() uuid.UUID {
//	job, err := c.Scheduler.NewJob(c.Duration, c.Task)
//	if err != nil {
//		fmt.Println("error adding job")
//		return uuid.UUID{}
//	}
//	return job.ID()
//}
//
