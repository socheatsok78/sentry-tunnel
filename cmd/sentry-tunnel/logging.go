package main

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var Default log.Logger

func init() {
	Default = log.NewLogfmtLogger(os.Stdout)
	Default = log.With(Default, "ts", log.DefaultTimestampUTC)
	Default = log.With(Default, "caller", log.DefaultCaller)
}

func SetLevel(option string) {
	switch option {
	case "debug":
		Default = level.NewFilter(Default, level.AllowDebug())
	case "info":
		Default = level.NewFilter(Default, level.AllowInfo())
	case "warn":
		Default = level.NewFilter(Default, level.AllowWarn())
	case "error":
		Default = level.NewFilter(Default, level.AllowError())
	case "none":
		Default = level.NewFilter(Default, level.AllowNone())
	default:
		Default = level.NewFilter(Default, level.AllowInfo())
	}
}
