package di

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	pw "wren-time-syncer/progress"
	"wren-time-syncer/renderer"
	"wren-time-syncer/repository"
	"wren-time-syncer/syncer"
	"wren-time-syncer/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/samber/do"
	"github.com/spf13/viper"
)

func provideConfig(path string) func(*do.Injector) (*utils.Config, error) {
	return func(i *do.Injector) (*utils.Config, error) {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}

		config := &utils.Config{}
		return config, viper.Unmarshal(config)
	}
}

func provideJiraClient(i *do.Injector) (*jira.Client, error) {
	config := do.MustInvoke[*utils.Config](i)
	tp := jira.BasicAuthTransport{
		Username: config.JiraUsername,
		Password: config.JiraToken,
	}
	client := tp.Client()
	timeout, err := strconv.Atoi(config.RequestTimeout)
	if err != nil {
		return nil, err
	}
	client.Timeout = time.Duration(timeout) * time.Millisecond
	jiraClient, err := jira.NewClient(client, config.JiraHost)
	if err != nil {
		return nil, err
	}
	return jiraClient, nil
}

func provideSyncer(i *do.Injector) (syncer.Syncer, error) {
	r := do.MustInvoke[repository.WorklogRepository](i)
	config := do.MustInvoke[*utils.Config](i)
	writer := do.MustInvoke[pw.Writer](i)
	return syncer.NewJiraSyncer(r, config, writer), nil
}

func provideWriter(_ *do.Injector) (io.Writer, error) {
	return os.Stdout, nil
}

func provideRenderer(i *do.Injector) (renderer.Writer, error) {
	wr := do.MustInvoke[io.Writer](i)
	return renderer.NewConsoleWriter(wr), nil
}

func provideProgressWriter(i *do.Injector) (pw.Writer, error) {
	w := do.MustInvoke[io.Writer](i)
	return pw.NewWriter(w), nil
}

func provideIssueRepository(i *do.Injector) (repository.IssueRepository, error) {
	client := do.MustInvoke[*jira.Client](i)
	config := do.MustInvoke[*utils.Config](i)
	writer := do.MustInvoke[pw.Writer](i)
	timeout, err := strconv.Atoi(config.RequestTimeout)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	return repository.NewJiraIssueRepository(client, httpClient, config, writer), nil
}

func provideWorklogRepository(i *do.Injector) (repository.WorklogRepository, error) {
	client := do.MustInvoke[*jira.Client](i)
	return repository.NewWorklogRepository(client), nil
}

func ProvideDeps(configPath string) error {
	do.Provide(nil, provideConfig(configPath))
	do.Provide(nil, provideJiraClient)
	do.Provide(nil, provideRenderer)
	do.Provide(nil, provideProgressWriter)
	do.Provide(nil, provideSyncer)
	do.Provide(nil, provideIssueRepository)
	do.Provide(nil, provideWorklogRepository)
	do.Provide(nil, provideWriter)
	return nil
}
