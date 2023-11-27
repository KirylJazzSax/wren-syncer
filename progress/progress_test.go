package progress

import (
	"os"
	"testing"

	mockprogresswriter "wren-time-syncer/mocks/go-pretty/progress"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)

	pw := mockprogresswriter.NewMockWriter(ctrl)

	prog := NewWriter(os.Stdout)

	pw.EXPECT().Render().Times(1)
	pw.EXPECT().AppendTracker(gomock.Any()).Times(1)

	err := prog.Start(20, "Tracker")
	assert.NotNil(t, prog.tr)
	assert.NoError(t, err)
}

func TestInc(t *testing.T) {
	ctrl := gomock.NewController(t)

	pw := mockprogresswriter.NewMockWriter(ctrl)

	prog := NewWriter(os.Stdout)

	pw.EXPECT().Render().Times(1)
	pw.EXPECT().AppendTracker(gomock.Any()).Times(1)

	length := 2
	err := prog.Start(int64(length), "Tracker")
	assert.NoError(t, err)

	err = prog.Inc(1)
	assert.NoError(t, err)

	pw.EXPECT().Stop().Times(1)

	err = prog.Inc(1)
	assert.NoError(t, err)
}

func TestStop(t *testing.T) {
	ctrl := gomock.NewController(t)

	pw := mockprogresswriter.NewMockWriter(ctrl)
	pw.EXPECT().Stop().Times(1)

	prog := NewWriter(os.Stdout)
	err := prog.Stop()
	assert.NoError(t, err)
}
