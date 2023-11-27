package progress

import (
	"io"
	"time"

	"wren-time-syncer/utils"

	p "github.com/jedib0t/go-pretty/v6/progress"
)

type PWriter struct {
	writer p.Writer
	tr     *p.Tracker
}

func (pw *PWriter) Start(l int64, name string) error {
	pw.tr = utils.NewTracker(name, l)
	pw.writer.AppendTracker(pw.tr)
	go pw.writer.Render()
	return nil
}

func (pw *PWriter) Inc(n int64) error {
	pw.tr.Increment(n)
	if pw.tr.IsDone() {
		return pw.Stop()
	}
	return nil
}

func (pw *PWriter) Stop() error {
	pw.writer.Stop()
	return nil
}

func NewWriter(writer io.Writer) *PWriter {
	prog := &PWriter{
		writer: newPW(writer),
	}
	return prog
}

func newPW(writer io.Writer) p.Writer {
	w := p.NewWriter()
	w.SetStyle(p.StyleCircle)
	w.SetTrackerPosition(p.PositionRight)
	w.SetUpdateFrequency(time.Millisecond * 100)
	w.SetOutputWriter(writer)
	w.Style().Colors = p.StyleColorsExample
	w.Style().Options.PercentFormat = "%4.1f%%"
	return w
}
