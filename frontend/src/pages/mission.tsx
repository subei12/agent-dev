import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { useState } from "react";

import { PageShell } from "../components/page-shell";
import {
  claimTask,
  createArchive,
  createDiscussionSession,
  createDocument,
  getDiscussionSessions,
  getMission,
  getMissionDocuments,
  getMissionRuntimes,
  getTaskBoard,
  requestReviewCheckpoint,
  sendTaskToAdmin
} from "../lib/api";
import { DiscussionSessionPanel } from "../features/discussions/discussion-session-panel";
import { DocumentList } from "../features/documents/document-list";
import { MissionActionRail } from "../features/missions/mission-action-rail";
import { MissionOverview } from "../features/missions/mission-overview";
import { TaskBoard } from "../features/tasks/task-board";
import { AgentRuntimeBoard } from "../features/runtime/agent-runtime-board";
import { RuntimeEventTimeline } from "../features/runtime/runtime-event-timeline";

const projectId = "proj_1";

export function MissionPage() {
  const { missionId = "mission_1" } = useParams();
  const queryClient = useQueryClient();
  const [taskActionMessage, setTaskActionMessage] = useState("");
  const [archiveStatus, setArchiveStatus] = useState("not archived");

  const missionQuery = useQuery({
    queryKey: ["mission", missionId],
    queryFn: () => getMission(projectId, missionId)
  });
  const documentsQuery = useQuery({
    queryKey: ["mission-documents", missionId],
    queryFn: () => getMissionDocuments(projectId, missionId)
  });
  const sessionsQuery = useQuery({
    queryKey: ["mission-discussions", missionId],
    queryFn: () => getDiscussionSessions(projectId, missionId)
  });
  const runtimeQuery = useQuery({
    queryKey: ["mission-runtimes", missionId],
    queryFn: () => getMissionRuntimes(projectId, missionId)
  });
  const taskBoardQuery = useQuery({
    queryKey: ["mission-task-board", missionId],
    queryFn: () => getTaskBoard(projectId, missionId)
  });

  const discussionMutation = useMutation({
    mutationFn: (topic: string) => createDiscussionSession(projectId, missionId, topic),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["mission-discussions", missionId] });
    }
  });

  const documentMutation = useMutation({
    mutationFn: ({ title, kind }: { title: string; kind: string }) => createDocument(projectId, missionId, title, kind),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["mission-documents", missionId] });
    }
  });

  const claimMutation = useMutation({
    mutationFn: (taskId: string) => claimTask(projectId, missionId, taskId),
    onSuccess: async () => {
      setTaskActionMessage("Task claimed");
      await queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
    }
  });

  const handoffMutation = useMutation({
    mutationFn: (taskId: string) => sendTaskToAdmin(projectId, missionId, taskId),
    onSuccess: async () => {
      setTaskActionMessage("Task handed to admin");
      await queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
    }
  });

  const checkpointMutation = useMutation({
    mutationFn: (taskId: string) => requestReviewCheckpoint(projectId, missionId, taskId),
    onSuccess: () => {
      setTaskActionMessage("Checkpoint requested");
    }
  });

  const archiveMutation = useMutation({
    mutationFn: () => createArchive(projectId, missionId),
    onSuccess: (result) => {
      setArchiveStatus(result.archive.status);
    }
  });

  return (
    <PageShell
      eyebrow="Mission Workspace"
      title="Documents, task lanes, and runtime telemetry"
      description="This workspace keeps the formal artifacts and live execution view on one surface so the admin agent can judge completion without leaving the page."
      aside={
        <div className="hero-stat">
          <span>Workspace</span>
          <strong>{missionId}</strong>
        </div>
      }
    >
      <MissionOverview mission={missionQuery.data ?? null} />
      <MissionActionRail
        archiveStatus={archiveStatus}
        isArchiving={archiveMutation.isPending}
        onArchive={() => archiveMutation.mutate()}
      />
      <DocumentList
        documents={documentsQuery.data ?? []}
        isSubmitting={documentMutation.isPending}
        onCreateDocument={(title, kind) => documentMutation.mutate({ title, kind })}
      />
      <DiscussionSessionPanel
        sessions={sessionsQuery.data ?? []}
        isSubmitting={discussionMutation.isPending}
        onCreateSession={(topic) => discussionMutation.mutate(topic)}
      />
      <TaskBoard
        tasks={taskBoardQuery.data?.items}
        taskActionMessage={taskActionMessage}
        onClaimTask={(taskId) => claimMutation.mutate(taskId)}
        onSendToAdmin={(taskId) => handoffMutation.mutate(taskId)}
        onRequestReview={(taskId) => checkpointMutation.mutate(taskId)}
      />
      <AgentRuntimeBoard runtimes={runtimeQuery.data ?? []} />
      <RuntimeEventTimeline events={[]} />
    </PageShell>
  );
}
