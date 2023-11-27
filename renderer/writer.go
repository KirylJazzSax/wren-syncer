package renderer

import (
	"context"

	"wren-time-syncer/repository"
)

type Writer interface {
	Render(context.Context, []repository.Issue) error
}
