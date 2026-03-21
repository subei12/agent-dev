package app

import (
	"context"
	"net/http"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformdb "github.com/your-org/agent-platform/internal/db"
	"github.com/your-org/agent-platform/internal/discussion"
	"github.com/your-org/agent-platform/internal/document"
	platformhttp "github.com/your-org/agent-platform/internal/http"
	"github.com/your-org/agent-platform/internal/mission"
	"github.com/your-org/agent-platform/internal/task"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIApp struct {
	server *http.Server
	pool   *pgxpool.Pool
}

func NewAPI(cfg platformconfig.Config) (*APIApp, error) {
	pool, err := platformdb.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	missionHandler := mission.NewHandler(
		mission.NewService(
			mission.NewRepository(pool),
		),
	)
	documentHandler := document.NewHandler(
		document.NewService(
			document.NewRepository(pool),
		),
	)
	discussionHandler := discussion.NewHandler(
		discussion.NewService(
			discussion.NewRepository(pool),
		),
	)
	taskHandler := task.NewHandler(
		task.NewService(
			task.NewRepository(pool),
		),
	)

	return &APIApp{
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: platformhttp.NewRouter(missionHandler, documentHandler, discussionHandler, taskHandler),
		},
		pool: pool,
	}, nil
}

func (a *APIApp) Run(_ context.Context) error {
	return a.server.ListenAndServe()
}
