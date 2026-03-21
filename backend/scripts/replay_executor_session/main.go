package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformdb "github.com/your-org/agent-platform/internal/db"
	"github.com/your-org/agent-platform/internal/runtime"
)

func main() {
	var (
		missionID         = flag.String("mission-id", "mission_1", "mission id")
		taskItemID        = flag.String("task-item-id", "task_runtime", "task item id")
		agentID           = flag.String("agent-id", "agent_backend", "agent id")
		executorProfileID = flag.String("executor-profile-id", "exec_backend", "executor profile id")
		backendName       = flag.String("backend", "codex_cli", "backend name")
	)
	flag.Parse()

	ctx := context.Background()
	cfg := platformconfig.MustLoad()

	pool, err := platformdb.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	service := runtime.NewService(runtime.NewRepository(pool))

	session, err := service.StartSession(ctx, runtime.StartSessionCmd{
		MissionID:         *missionID,
		TaskItemID:        *taskItemID,
		AgentID:           *agentID,
		ExecutorProfileID: *executorProfileID,
		Backend:           *backendName,
		Status:            "starting",
	})
	if err != nil {
		log.Fatal(err)
	}

	events := []runtime.AppendEventCmd{
		{
			MissionID:         *missionID,
			AgentID:           *agentID,
			ExecutorSessionID: session.ID,
			TaskItemID:        *taskItemID,
			Level:             "info",
			Category:          "session",
			Type:              "session_started",
			Title:             "Session started",
			Summary:           "seeded replay launched",
		},
		{
			MissionID:         *missionID,
			AgentID:           *agentID,
			ExecutorSessionID: session.ID,
			TaskItemID:        *taskItemID,
			Level:             "info",
			Category:          "model",
			Type:              "message_created",
			Title:             "Model output",
			Summary:           "assistant produced a seeded runtime summary",
		},
	}

	for _, event := range events {
		if err := service.AppendEvent(ctx, event); err != nil {
			log.Fatal(err)
		}
	}

	entries := []runtime.AppendTranscriptEntryCmd{
		{
			ExecutorSessionID: session.ID,
			Role:              "assistant",
			EntryType:         "message",
			ContentText:       "Drafted runtime transcript for admin review.",
		},
		{
			ExecutorSessionID: session.ID,
			Role:              "assistant",
			EntryType:         "message",
			ContentText:       "token=abc123 should be redacted in transcript storage.",
		},
	}

	for _, entry := range entries {
		if err := service.AppendTranscriptEntry(ctx, entry); err != nil {
			log.Fatal(err)
		}
	}

	if err := service.SealSession(ctx, session.ID); err != nil {
		log.Fatal(err)
	}

	if _, err := service.GetTranscriptView(ctx, session.ID, "demo-admin", "redacted"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("replayed executor session: %s\n", session.ID)
}
