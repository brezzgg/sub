package badger

import (
	"strings"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/repo/repo"
	"github.com/dgraph-io/badger/v2"
)

type Logger struct{}

// Debugf implements [badger.Logger].
func (l *Logger) Debugf(f string, a ...any) {
	repo.LogLevel.Log(
		lg.GlobalLogger,
		lg.F(strings.TrimSpace(f), a...),
		lg.C{"level": "debug"},
	)
}

// Errorf implements [badger.Logger].
func (l *Logger) Errorf(f string, a ...any) {
	repo.LogLevel.Log(
		lg.GlobalLogger,
		lg.F(strings.TrimSpace(f), a...),
		lg.C{"level": "error"},
	)
}

// Infof implements [badger.Logger].
func (l *Logger) Infof(f string, a ...any) {
	repo.LogLevel.Log(
		lg.GlobalLogger,
		lg.F(strings.TrimSpace(f), a...),
		lg.C{"level": "info"},
	)
}

// Warningf implements [badger.Logger].
func (l *Logger) Warningf(f string, a ...any) {
	repo.LogLevel.Log(
		lg.GlobalLogger,
		lg.F(strings.TrimSpace(f), a...),
		lg.C{"level": "warn"},
	)
}

var _ badger.Logger = (*Logger)(nil)
