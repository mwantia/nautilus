package log

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

var Default hclog.Logger

type HighlightFormatter struct {
	NoColor bool `json:"no_color"`
}

type Logger struct {
	hclog.Logger
}

func Setup(level string, noColor bool) error {
	logLevel := ParseLevel(level)

	Default = hclog.New(&hclog.LoggerOptions{
		Name:        "nautilus",
		Level:       logLevel,
		Output:      os.Stdout,
		JSONFormat:  false,
		Color:       hclog.AutoColor,
		TimeFormat:  "02.01.2006 15:04:05",
		DisableTime: false,
	})

	log.SetOutput(io.Discard)
	hclog.SetDefault(Default)

	return nil
}

func ParseLevel(level string) hclog.Level {
	switch strings.ToUpper(level) {
	case "TRACE":
		return hclog.Trace
	case "DEBUG":
		return hclog.Debug
	case "INFO":
		return hclog.Info
	case "WARN":
		return hclog.Warn
	case "ERROR":
		return hclog.Error
	default:
		return hclog.Info
	}
}

func NewLogger(name string) *Logger {
	if Default != nil {
		return &Logger{
			Logger: Default.Named(name),
		}
	}

	return &Logger{
		Logger: hclog.Default(),
	}
}
