package cmd

import (
	"context"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"wren-time-syncer/internal/di"
	mockrenderer "wren-time-syncer/mocks/renderer"
	"wren-time-syncer/renderer"
	"wren-time-syncer/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/golang/mock/gomock"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestSyncCmd(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()

	utils.MustVoid(di.ProvideDeps("."))

	mux := utils.MakeTestServeMux(2, false, 200)
	ts := httptest.NewServer(mux)

	do.Override(nil, func(i *do.Injector) (*utils.Config, error) {
		config := &utils.Config{
			JiraToken:      "token",
			LinkAuthHeader: "header",
			JiraHost:       ts.URL,
			LinkHost:       ts.URL + "/",
			RequestTimeout: "1000",
		}

		return config, nil
	})

	do.Override(nil, func(i *do.Injector) (*jira.Client, error) {
		config := do.MustInvoke[*utils.Config](i)

		timeout, err := strconv.Atoi(config.RequestTimeout)
		if err != nil {
			return nil, err
		}
		client := ts.Client()
		client.Timeout = time.Duration(timeout) * time.Millisecond
		jiraClient, err := jira.NewClient(client, config.JiraHost)
		if err != nil {
			return nil, err
		}
		return jiraClient, nil
	})

	do.Override(nil, func(i *do.Injector) (renderer.Writer, error) {
		w := mockrenderer.NewMockWriter(ctrl)
		w.EXPECT().Render(gomock.Any(), gomock.Any()).Times(1).Return(nil)
		return w, nil
	})

	h, err := newHandler("sync", ctx)
	assert.NoError(t, err)

	cmd := newSyncCmd(h)
	assert.NoError(t, cmd.Execute())
}

func TestExecute(t *testing.T) {
	err := Execute(context.TODO())
	assert.NoError(t, err)
}
