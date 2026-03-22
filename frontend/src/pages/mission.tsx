import { useMutation, useQueries, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";

import { PageShell } from "../components/page-shell";
import {
  claimTask,
  createArchive,
  createDiscussionSession,
  createDocument,
  createDocumentVersion,
  getDiscussionSessions,
  getDocumentVersions,
  getMissionApprovals,
  getMissionArchive,
  getMission,
  getMissionDocuments,
  getMissionRuntimes,
  getTaskBoard,
  adoptDocumentVersion,
  requestReviewCheckpoint,
  sendTaskToAdmin
} from "../lib/api";
import { subscribeToChannel } from "../lib/sse";
import { DiscussionSessionPanel } from "../features/discussions/discussion-session-panel";
import { DocumentList } from "../features/documents/document-list";
import { MissionActionRail } from "../features/missions/mission-action-rail";
import { MissionApprovalPanel } from "../features/missions/mission-approval-panel";
import { MissionArchivePanel } from "../features/missions/mission-archive-panel";
import { MissionOverview } from "../features/missions/mission-overview";
import { TaskBoard } from "../features/tasks/task-board";
import { AgentRuntimeBoard } from "../features/runtime/agent-runtime-board";
import { RuntimeEventTimeline } from "../features/runtime/runtime-event-timeline";

const projectId = "proj_1";

/**
 * MissionPage 渲染当前路由对应的页面级工作区。
 */
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
  const approvalQuery = useQuery({
    queryKey: ["mission-approvals", missionId],
    queryFn: () => getMissionApprovals(projectId, missionId)
  });
  const archiveQuery = useQuery({
    queryKey: ["mission-archive", missionId],
    queryFn: () => getMissionArchive(projectId, missionId)
  });
  const taskBoardQuery = useQuery({
    queryKey: ["mission-task-board", missionId],
    queryFn: () => getTaskBoard(projectId, missionId)
  });
  const versionQueries = useQueries({
    queries: (documentsQuery.data ?? []).map((document) => ({
      queryKey: ["document-versions", missionId, document.id],
      queryFn: () => getDocumentVersions(projectId, missionId, document.id)
    }))
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
  const versionMutation = useMutation({
    mutationFn: ({ documentId, contentText }: { documentId: string; contentText: string }) =>
      createDocumentVersion(projectId, missionId, documentId, contentText),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({ queryKey: ["document-versions", missionId, variables.documentId] });
    }
  });
  const adoptMutation = useMutation({
    mutationFn: ({ documentId, versionId }: { documentId: string; versionId: string }) =>
      adoptDocumentVersion(projectId, missionId, documentId, versionId),
    onSuccess: async (_, variables) => {
      await queryClient.invalidateQueries({ queryKey: ["document-versions", missionId, variables.documentId] });
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
      void queryClient.invalidateQueries({ queryKey: ["mission-archive", missionId] });
    }
  });
  const versionsByDocument = Object.fromEntries(
    (documentsQuery.data ?? []).map((document, index) => [document.id, versionQueries[index]?.data ?? []])
  );

  useEffect(() => {
    return subscribeToChannel(projectId, `mission:${missionId}`, () => {
      void queryClient.invalidateQueries({ queryKey: ["mission-runtimes", missionId] });
      void queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
      void queryClient.invalidateQueries({ queryKey: ["mission-approvals", missionId] });
      void queryClient.invalidateQueries({ queryKey: ["mission-archive", missionId] });
    });
  }, [missionId, queryClient]);

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
      <MissionApprovalPanel approvals={approvalQuery.data ?? []} />
      <MissionArchivePanel archive={archiveQuery.data ?? null} />
      <DocumentList
        documents={documentsQuery.data ?? []}
        isSubmitting={documentMutation.isPending || versionMutation.isPending || adoptMutation.isPending}
        versionsByDocument={versionsByDocument}
        onCreateDocument={(title, kind) => documentMutation.mutate({ title, kind })}
        onCreateVersion={(documentId, contentText) => versionMutation.mutate({ documentId, contentText })}
        onAdoptVersion={(documentId, versionId) => adoptMutation.mutate({ documentId, versionId })}
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
