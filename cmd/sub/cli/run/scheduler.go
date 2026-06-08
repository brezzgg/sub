package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/brezzgg/go-packages/lg"
)

type task struct {
	run  func() error
	stop func() error
}

type scheduler struct {
	tasks []*task
	errCh chan error
}

func (s *scheduler) Add(runFn, stopFn func() error) {
	s.tasks = append(s.tasks, &task{
		run:  runFn,
		stop: stopFn,
	})
}

func (s *scheduler) Run() error {
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

func (s *scheduler) wait() error {
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

func (s *scheduler) Stop() {
	for _, t := range s.tasks {
		_ = t.stop()
	}
}
