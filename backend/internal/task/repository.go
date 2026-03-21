package task

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateBoard(ctx context.Context, missionID, title string) (TaskBoard, error) {
	row, err := sqlc.New(r.pool).CreateTaskBoard(ctx, sqlc.CreateTaskBoardParams{
		ID:        uuid.NewString(),
		MissionID: missionID,
		Title:     title,
	})
	if err != nil {
		return TaskBoard{}, err
	}
	return boardFromRow(row), nil
}

func (r *Repository) CreateTaskItem(ctx context.Context, boardID string, cmd CreateTaskItemCmd) (TaskItem, error) {
	row, err := sqlc.New(r.pool).CreateTaskItem(ctx, sqlc.CreateTaskItemParams{
		ID:                        uuid.NewString(),
		BoardID:                   boardID,
		Title:                     cmd.Title,
		Type:                      cmd.Type,
		Status:                    "todo",
		AssignedAgentID:           textValue(cmd.AssignedAgentID),
		UpstreamTaskIdsJson:       jsonValue(cmd.UpstreamTaskIDs, []byte("[]")),
		DownstreamTaskIdsJson:     jsonValue(cmd.DownstreamTaskIDs, []byte("[]")),
		InputDocumentVersionIdsJson: jsonValue(cmd.InputDocumentVersionIDs, []byte("[]")),
		InputRepoCandidateIdsJson:   jsonValue(cmd.InputRepoCandidateIDs, []byte("[]")),
		DefinitionOfDoneJson:        jsonValue(cmd.DefinitionOfDone, []byte("{}")),
	})
	if err != nil {
		return TaskItem{}, err
	}
	return itemFromRow(row), nil
}

func (r *Repository) GetBoard(ctx context.Context, missionID string) (TaskBoard, error) {
	row, err := sqlc.New(r.pool).GetLatestTaskBoardByMission(ctx, missionID)
	if err != nil {
		return TaskBoard{}, err
	}
	return boardFromRow(row), nil
}

func (r *Repository) ListTaskItems(ctx context.Context, boardID string) ([]TaskItem, error) {
	rows, err := sqlc.New(r.pool).ListTaskItemsByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}
	items := make([]TaskItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, itemFromRow(row))
	}
	return items, nil
}

func (r *Repository) GetTaskItem(ctx context.Context, taskItemID string) (TaskItem, error) {
	row, err := sqlc.New(r.pool).GetTaskItem(ctx, taskItemID)
	if err != nil {
		return TaskItem{}, err
	}
	return itemFromRow(row), nil
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, taskItemID, status string) (TaskItem, error) {
	row, err := sqlc.New(r.pool).UpdateTaskItemStatus(ctx, sqlc.UpdateTaskItemStatusParams{
		ID:     taskItemID,
		Status: status,
	})
	if err != nil {
		return TaskItem{}, err
	}
	return itemFromRow(row), nil
}

func (r *Repository) CreateClaim(ctx context.Context, cmd ClaimTaskCmd) (TaskClaim, error) {
	row, err := sqlc.New(r.pool).CreateTaskClaim(ctx, sqlc.CreateTaskClaimParams{
		ID:          uuid.NewString(),
		TaskItemID:  cmd.TaskItemID,
		AgentID:     cmd.AgentID,
		Status:      "active",
		ClaimReason: textValue(cmd.ClaimReason),
	})
	if err != nil {
		return TaskClaim{}, err
	}
	return claimFromCreateRow(row), nil
}

func (r *Repository) CreateHandoff(ctx context.Context, cmd CreateHandoffCmd) (TaskHandoff, error) {
	row, err := sqlc.New(r.pool).CreateTaskHandoff(ctx, sqlc.CreateTaskHandoffParams{
		ID:                         uuid.NewString(),
		TaskItemID:                 cmd.TaskItemID,
		FromAgentID:                cmd.FromAgentID,
		ToAgentID:                  textValue(cmd.ToAgentID),
		ToAdminAgent:               cmd.ToAdminAgent,
		Summary:                    cmd.Summary,
		OutputDocumentVersionIdsJson: jsonValue(cmd.OutputDocumentVersionIDs, []byte("[]")),
		OutputRepoCandidateIdsJson:   jsonValue(cmd.OutputRepoCandidateIDs, []byte("[]")),
		OutputArtifactVersionIdsJson: jsonValue(cmd.OutputArtifactVersionIDs, []byte("[]")),
		ValidationSummary:            textValue(cmd.ValidationSummary),
		RiskSummary:                  textValue(cmd.RiskSummary),
		RecommendedNextAction:        textValue(cmd.RecommendedNextAction),
		Status:                       "pending",
	})
	if err != nil {
		return TaskHandoff{}, err
	}
	return handoffFromRow(row), nil
}

func (r *Repository) CreateCheckpoint(ctx context.Context, cmd RequestCheckpointCmd) (ReviewCheckpoint, error) {
	row, err := sqlc.New(r.pool).CreateReviewCheckpoint(ctx, sqlc.CreateReviewCheckpointParams{
		ID:                        uuid.NewString(),
		MissionID:                 cmd.MissionID,
		TaskItemID:                textValue(cmd.TaskItemID),
		Kind:                      cmd.Kind,
		RequestedByAgentID:        cmd.RequestedByAgentID,
		AssignedAgentID:           textValue(cmd.AssignedAgentID),
		Status:                    "pending",
		Summary:                   textValue(cmd.Summary),
		LinkedDocumentVersionIdsJson: jsonValue(cmd.LinkedDocumentVersionIDs, []byte("[]")),
		LinkedRepoCandidateIdsJson:   jsonValue(cmd.LinkedRepoCandidateIDs, []byte("[]")),
	})
	if err != nil {
		return ReviewCheckpoint{}, err
	}
	return checkpointFromRow(row), nil
}

func boardFromRow(row sqlc.TaskBoard) TaskBoard {
	return TaskBoard{
		ID:        row.ID,
		MissionID: row.MissionID,
		Title:     row.Title,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func itemFromRow(row sqlc.TaskItem) TaskItem {
	return TaskItem{
		ID:                      row.ID,
		BoardID:                 row.BoardID,
		Title:                   row.Title,
		Type:                    row.Type,
		Status:                  row.Status,
		AssignedAgentID:         stringValue(row.AssignedAgentID),
		UpstreamTaskIDs:         json.RawMessage(row.UpstreamTaskIdsJson),
		DownstreamTaskIDs:       json.RawMessage(row.DownstreamTaskIdsJson),
		InputDocumentVersionIDs: json.RawMessage(row.InputDocumentVersionIdsJson),
		InputRepoCandidateIDs:   json.RawMessage(row.InputRepoCandidateIdsJson),
		DefinitionOfDone:        json.RawMessage(row.DefinitionOfDoneJson),
		CreatedAt:               row.CreatedAt.Time,
		UpdatedAt:               row.UpdatedAt.Time,
	}
}

func claimFromRow(row sqlc.TaskClaim) TaskClaim {
	return TaskClaim{
		ID:          row.ID,
		TaskItemID:  row.TaskItemID,
		AgentID:     row.AgentID,
		Status:      row.Status,
		ClaimReason: stringValue(row.ClaimReason),
		CreatedAt:   row.CreatedAt.Time,
	}
}

func claimFromCreateRow(row sqlc.CreateTaskClaimRow) TaskClaim {
	return TaskClaim{
		ID:          row.ID,
		TaskItemID:  row.TaskItemID,
		AgentID:     row.AgentID,
		Status:      row.Status,
		ClaimReason: stringValue(row.ClaimReason),
		CreatedAt:   row.CreatedAt.Time,
	}
}

func handoffFromRow(row sqlc.TaskHandoff) TaskHandoff {
	return TaskHandoff{
		ID:                       row.ID,
		TaskItemID:               row.TaskItemID,
		FromAgentID:              row.FromAgentID,
		ToAgentID:                stringValue(row.ToAgentID),
		ToAdminAgent:             row.ToAdminAgent,
		Summary:                  row.Summary,
		OutputDocumentVersionIDs: json.RawMessage(row.OutputDocumentVersionIdsJson),
		OutputRepoCandidateIDs:   json.RawMessage(row.OutputRepoCandidateIdsJson),
		OutputArtifactVersionIDs: json.RawMessage(row.OutputArtifactVersionIdsJson),
		ValidationSummary:        stringValue(row.ValidationSummary),
		RiskSummary:              stringValue(row.RiskSummary),
		RecommendedNextAction:    stringValue(row.RecommendedNextAction),
		Status:                   row.Status,
		CreatedAt:                row.CreatedAt.Time,
	}
}

func checkpointFromRow(row sqlc.ReviewCheckpoint) ReviewCheckpoint {
	return ReviewCheckpoint{
		ID:                       row.ID,
		MissionID:                row.MissionID,
		TaskItemID:               stringValue(row.TaskItemID),
		Kind:                     row.Kind,
		RequestedByAgentID:       row.RequestedByAgentID,
		AssignedAgentID:          stringValue(row.AssignedAgentID),
		Status:                   row.Status,
		Summary:                  stringValue(row.Summary),
		LinkedDocumentVersionIDs: json.RawMessage(row.LinkedDocumentVersionIdsJson),
		LinkedRepoCandidateIDs:   json.RawMessage(row.LinkedRepoCandidateIdsJson),
		CreatedAt:                row.CreatedAt.Time,
	}
}
