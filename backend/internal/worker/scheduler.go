package worker

import "time"

type Scheduler struct {
	pollInterval time.Duration
}

func NewScheduler(pollInterval time.Duration) *Scheduler {
	return &Scheduler{pollInterval: pollInterval}
}

func (s *Scheduler) PollInterval() time.Duration {
	return s.pollInterval
}
