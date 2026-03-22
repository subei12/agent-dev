package agentcfg

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type AgentConfig struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Enabled           bool     `json:"enabled"`
	ExecutorProfileID string   `json:"executorProfileId"`
	ExecutorType      string   `json:"executorType"`
	Command           string   `json:"command"`
	Args              []string `json:"args"`
	TimeoutSec        int32    `json:"timeoutSec"`
	MaxConcurrency    int32    `json:"maxConcurrency"`
	AllowRepoRead     bool     `json:"allowRepoRead"`
	AllowRepoWrite    bool     `json:"allowRepoWrite"`
	AllowNetwork      bool     `json:"allowNetwork"`
}

type UpdateAgentConfigCmd struct {
	AgentID           string
	Name              string
	Enabled           bool
	ExecutorType      string
	Command           string
	Args              []string
	TimeoutSec        int32
	MaxConcurrency    int32
	AllowRepoRead     bool
	AllowRepoWrite    bool
	AllowNetwork      bool
}

type Service interface {
	List(context.Context, string) ([]AgentConfig, error)
	Update(context.Context, UpdateAgentConfigCmd) (AgentConfig, error)
}

type service struct {
	queries *sqlc.Queries
}

// NewService 创建并返回 Agent 配置服务。
func NewService(pool *pgxpool.Pool) Service {
	return &service{queries: sqlc.New(pool)}
}

// List 返回项目下所有 Agent 的配置。
func (s *service) List(ctx context.Context, projectID string) ([]AgentConfig, error) {
	agents, err := s.queries.ListAgentsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	profiles, err := s.queries.ListExecutorProfilesByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	profileByID := make(map[string]sqlc.ExecutorProfile, len(profiles))
	for _, profile := range profiles {
		profileByID[profile.ID] = profile
	}

	items := make([]AgentConfig, 0, len(agents))
	for _, agent := range agents {
		profile := profileByID[agent.ExecutorProfileID]
		items = append(items, mapAgent(agent, profile))
	}
	return items, nil
}

// Update 更新 Agent 及其执行器配置。
func (s *service) Update(ctx context.Context, cmd UpdateAgentConfigCmd) (AgentConfig, error) {
	agent, err := s.queries.UpdateAgent(ctx, sqlc.UpdateAgentParams{
		ID:      cmd.AgentID,
		Name:    cmd.Name,
		Enabled: cmd.Enabled,
	})
	if err != nil {
		return AgentConfig{}, err
	}

	profile, err := s.queries.UpdateExecutorProfile(ctx, sqlc.UpdateExecutorProfileParams{
		ID:              agent.ExecutorProfileID,
		Type:            cmd.ExecutorType,
		Command:         cmd.Command,
		ArgsJson:        jsonBytes(cmd.Args),
		TimeoutSec:      cmd.TimeoutSec,
		MaxConcurrency:  cmd.MaxConcurrency,
		AllowRepoRead:   cmd.AllowRepoRead,
		AllowRepoWrite:  cmd.AllowRepoWrite,
		AllowNetwork:    cmd.AllowNetwork,
	})
	if err != nil {
		return AgentConfig{}, err
	}

	return mapAgent(agent, profile), nil
}

// mapAgent 把数据库模型映射成前端需要的 Agent 配置视图。
func mapAgent(agent sqlc.Agent, profile sqlc.ExecutorProfile) AgentConfig {
	return AgentConfig{
		ID:                agent.ID,
		Name:              agent.Name,
		Enabled:           agent.Enabled,
		ExecutorProfileID: agent.ExecutorProfileID,
		ExecutorType:      profile.Type,
		Command:           profile.Command,
		Args:              jsonStrings(profile.ArgsJson),
		TimeoutSec:        profile.TimeoutSec,
		MaxConcurrency:    profile.MaxConcurrency,
		AllowRepoRead:     profile.AllowRepoRead,
		AllowRepoWrite:    profile.AllowRepoWrite,
		AllowNetwork:      profile.AllowNetwork,
	}
}

func jsonStrings(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var items []string
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil
	}
	return items
}

func jsonBytes(items []string) []byte {
	body, err := json.Marshal(items)
	if err != nil {
		return []byte("[]")
	}
	return body
}
