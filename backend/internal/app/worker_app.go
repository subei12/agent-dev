package app

import (
	"context"
	"time"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformdb "github.com/your-org/agent-platform/internal/db"
	"github.com/your-org/agent-platform/internal/executor"
	"github.com/your-org/agent-platform/internal/runtime"
	"github.com/your-org/agent-platform/internal/worker"
)

type WorkerApp struct {
	cfg       platformconfig.Config
	scheduler *worker.Scheduler
}

func NewWorker(cfg platformconfig.Config) (*WorkerApp, error) {
	pool, err := platformdb.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	runtimeService := runtime.NewService(runtime.NewRepository(pool))
	manager := executor.NewManager(
		executor.NewBaseAdapter("codex_cli"),
		runtimeService,
		executor.NewPtyRunner(),
	)
	repo := worker.NewRepository(pool)
	runService := worker.NewRunService(repo, worker.NewExecutorLauncher(manager))
	dispatcher := worker.NewDispatcher(runService)

	return &WorkerApp{
		cfg:       cfg,
		scheduler: worker.NewScheduler(repo, dispatcher),
	}, nil
}

func (w *WorkerApp) Run(ctx context.Context) error {
	return w.scheduler.Run(ctx, 5*time.Second)
}
