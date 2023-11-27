package utils

import "github.com/jedib0t/go-pretty/v6/progress"

func NewTracker(m string, l int64) *progress.Tracker {
	units := &progress.UnitsDefault
	return &progress.Tracker{Message: m, Total: l, Units: *units}
}
