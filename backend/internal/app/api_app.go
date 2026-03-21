package app

import (
	"context"
	"net/http"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformdb "github.com/your-org/agent-platform/internal/db"
	platformhttp "github.com/your-org/agent-platform/internal/http"
	"github.com/your-org/agent-platform/internal/mission"
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

	return &APIApp{
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: platformhttp.NewRouter(missionHandler),
		},
		pool: pool,
	}, nil
}

func (a *APIApp) Run(_ context.Context) error {
	return a.server.ListenAndServe()
}
