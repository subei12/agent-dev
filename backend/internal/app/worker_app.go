package app

import (
	"context"

	platformconfig "github.com/your-org/agent-platform/internal/config"
)

type WorkerApp struct {
	cfg platformconfig.Config
}

func NewWorker(cfg platformconfig.Config) *WorkerApp {
	return &WorkerApp{cfg: cfg}
}

func (w *WorkerApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}
