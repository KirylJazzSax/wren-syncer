package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"wren-time-syncer/cmd"
	"wren-time-syncer/internal/di"
	"wren-time-syncer/utils"
)

//go:generate mockgen -source=$GOPATH/pkg/mod/github.com/jedib0t/go-pretty/v6@v6.4.4/progress/writer.go -destination=$GOPATH/src/wren-syncer/mocks/go-pretty/progress/writer.go
//go:generate mockgen -source=$GOPATH/pkg/mod/github.com/jedib0t/go-pretty/v6@v6.4.4/table/writer.go -destination=$GOPATH/src/wren-syncer/mocks/go-pretty/table/writer.go
//go:generate mockgen -source=$GOPATH/src/wren-syncer/progress/progress.go -destination=$GOPATH/src/wren-syncer/mocks/progress/progress.go
//go:generate mockgen -source=$GOPATH/src/wren-syncer/renderer/writer.go -destination=$GOPATH/src/wren-syncer/mocks/renderer/writer.go
//go:generate mockgen -source=$GOPATH/src/wren-syncer/repository/repository.go -destination=$GOPATH/src/wren-syncer/mocks/repository/repository.go
//go:generate mockgen -source=$GOPATH/src/wren-syncer/syncer/syncer.go -destination=$GOPATH/src/wren-syncer/mocks/syncer/syncer.go
//go:generate mockgen -source=$GOPATH/src/wren-syncer/prompt/confirm_prompt_runner.go -destination=$GOPATH/src/wren-syncer/mocks/prompt/confirm_prompt_runner.go
//go:generate mockgen -source=$GOPATH/src/wren-syncer/prompt/select_prompt_runner.go -destination=$GOPATH/src/wren-syncer/mocks/prompt/select_prompt_runner.go
func main() {
	utils.MustVoid(di.ProvideDeps("."))
	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)
	signal.Notify(
		exit,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go func() {
		<-exit
		cancel()
		os.Exit(0)
	}()

	if err := cmd.Execute(ctx); err != nil {
		utils.LogError(err.Error())
		os.Exit(1)
	}
}
