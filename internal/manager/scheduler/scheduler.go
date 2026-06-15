package scheduler

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/brezzgg/go-packages/lg"
)

type Task struct {
	run  func() error
	stop func() error
}

type Scheduler struct {
	tasks []*Task
	errCh chan error
}

func (s *Scheduler) Add(runFn, stopFn func() error) {
	s.tasks = append(s.tasks, &Task{
		run:  runFn,
		stop: stopFn,
	})
}

func (s *Scheduler) Run() error {
	if s.errCh == nil {
		s.errCh = make(chan error, len(s.tasks))
	}
	for _, t := range s.tasks {
		go func() {
			s.errCh <- t.run()
		}()
	}
	return s.wait()
}

func (s *Scheduler) wait() error {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-s.errCh:
		s.Stop()
		return err
	case <-sigCh:
		lg.Info("graceful shutdown...", lg.Sync{})
		s.Stop()
		return nil
	}
}

func (s *Scheduler) Stop() {
	for _, t := range s.tasks {
		_ = t.stop()
	}
}
