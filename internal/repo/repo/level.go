package repo

import "github.com/brezzgg/go-packages/lg"

var LogLevel = lg.NewLogLevel(
	lg.ClrFgBoldYellow,
	"Repo",
	lg.LevelOptionCallerOnlyFunc,
)
