package cdflog

import (
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"os"
	"time"
	cdflog "upgrade/pkg/log"
)

type ZeroLog struct {
	zerolog.Logger
	level zerolog.Level
	w     zerolog.LevelWriter
}

// New creates a root logger with given output writer. If the output writer implements
// the LevelWriter interface, the WriteLevel method will be called instead of the Write
// one.
//
// Each logging operation makes a single call to the Writer's Write method. There is no
// guarantee on access serialization to the Writer. If your Writer is not thread safe,
// you may consider using sync wrapper.
func New(w io.Writer) ZeroLog {
	if w == nil {
		w = ioutil.Discard
	}
	lw, ok := w.(zerolog.LevelWriter)
	if !ok {
		lw = levelWriterAdapter{w}
	}
	return ZeroLog{w: lw, level: zerolog.TraceLevel}
}

type levelWriterAdapter struct {
	io.Writer
}

func (lw levelWriterAdapter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	return lw.Write(p)
}

func NewZeroLog(logfile *os.File, level cdflog.Level) *ZeroLog {
	var writer []io.Writer

	writer = append(writer, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	if logfile != nil {
		writer = append(writer, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	multi := zerolog.MultiLevelWriter(writer...)

	logger := New(multi)

	return &logger
}
