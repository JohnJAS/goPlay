package cdflog

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

const (
	LogLevelEnv = "LOG_LEVEL"
)

var (
	EnvLogLevel = os.Getenv(LogLevelEnv)
)

type WriterList []io.Writer

type FilteredWriter struct {
	zerolog.ConsoleWriter
	level zerolog.Level
}

func (w FilteredWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.level {
		return w.Write(p)
	}
	return len(p), nil
}

func NewZeroLog(file *os.File, level zerolog.Level) zerolog.Logger {

	writers := NewWriter(file)

	levelWriter := zerolog.MultiLevelWriter(writers...)

	logger := zerolog.New(levelWriter).With().Timestamp().Logger()

	if level == 0 {
		if EnvLogLevel != "" {
			lvl, _ := zerolog.ParseLevel(EnvLogLevel)
			zerolog.SetGlobalLevel(lvl)
		}
	}

	return logger
}

func NewWriter(file *os.File) (wl WriterList) {
	fw := FilteredWriter{
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		},
		zerolog.InfoLevel,
	}
	fw.FormatLevel = func(interface{}) string {
		return ""
	}
	fw.FormatTimestamp = func(interface{}) string {
		return ""
	}
	wl = append(wl, fw)

	if file != nil {
		wl = append(wl, zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		})
	}

	return
}
