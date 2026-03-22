package app

import (
	"context"
	"net/http"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformdb "github.com/your-org/agent-platform/internal/db"
	"github.com/your-org/agent-platform/internal/approval"
	"github.com/your-org/agent-platform/internal/agentcfg"
	"github.com/your-org/agent-platform/internal/authz"
	"github.com/your-org/agent-platform/internal/archive"
	"github.com/your-org/agent-platform/internal/discussion"
	"github.com/your-org/agent-platform/internal/document"
	platformhttp "github.com/your-org/agent-platform/internal/http"
	"github.com/your-org/agent-platform/internal/mission"
	"github.com/your-org/agent-platform/internal/runtime"
	"github.com/your-org/agent-platform/internal/sse"
	platformstorage "github.com/your-org/agent-platform/internal/storage"
	"github.com/your-org/agent-platform/internal/task"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIApp struct {
	server *http.Server
	pool   *pgxpool.Pool
}

// NewAPI 创建并返回对应的组件。
func NewAPI(cfg platformconfig.Config) (*APIApp, error) {
	pool, err := platformdb.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	hub := sse.NewHub()
	objectStore, err := platformstorage.NewObjectStore(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket)
	if err != nil {
		return nil, err
	}
	authorizer := authz.New(pool)

	missionHandler := mission.NewHandler(
		mission.NewService(
			mission.NewRepository(pool),
		),
		authorizer,
	)
	documentHandler := document.NewHandler(
		document.NewService(
			document.NewRepository(pool),
		),
		authorizer,
	)
	discussionHandler := discussion.NewHandler(
		discussion.NewService(
			discussion.NewRepository(pool),
		),
		authorizer,
	)
	taskHandler := task.NewHandler(
		task.NewService(
			task.NewRepository(pool),
		),
		authorizer,
	)
	runtimeHandler := runtime.NewHandler(
		runtime.NewServiceWithDeps(
			runtime.NewRepository(pool),
			hub,
			objectStore,
		),
		authorizer,
	)
	sseHandler := sse.NewHandler(hub)
	approvalHandler := approval.NewHandler(
		approval.NewService(pool),
	)
	agentConfigHandler := agentcfg.NewHandler(
		agentcfg.NewService(pool),
		authorizer,
	)
	archiveHandler := archive.NewHandler(
		archive.NewService(
			archive.NewRepository(pool),
			objectStore,
		),
		authorizer,
	)

	return &APIApp{
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: platformhttp.NewRouter(missionHandler, documentHandler, discussionHandler, taskHandler, runtimeHandler, approvalHandler, agentConfigHandler, archiveHandler, sseHandler),
		},
		pool: pool,
	}, nil
}

// Run 执行当前组件的主循环或工作流。
func (a *APIApp) Run(_ context.Context) error {
	return a.server.ListenAndServe()
}
