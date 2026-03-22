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
  createTask,
  getDiscussionSessions,
  getDocumentVersions,
  getAgentConfigs,
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
import { TaskComposer } from "../features/tasks/task-composer";
import { TaskDetailPanel } from "../features/tasks/task-detail-panel";
import { TaskBoard } from "../features/tasks/task-board";
import { AgentRuntimeBoard } from "../features/runtime/agent-runtime-board";
import { RuntimeEventTimeline } from "../features/runtime/runtime-event-timeline";
import { TaskRuntimePanel } from "../features/runtime/task-runtime-panel";
import { getRuntimeEvents, getRuntimeSession, getTranscript, getTranscriptAccessAudits } from "../lib/api";

const projectId = "proj_1";

/**
 * MissionPage 渲染当前路由对应的页面级工作区。
 */
export function MissionPage() {
  const { missionId = "mission_1" } = useParams();
  const queryClient = useQueryClient();
  const [taskActionMessage, setTaskActionMessage] = useState("");
  const [archiveStatus, setArchiveStatus] = useState("未归档");
  const [selectedTaskId, setSelectedTaskId] = useState<string | null>(null);
  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);

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
  const agentsQuery = useQuery({
    queryKey: ["agent-configs"],
    queryFn: () => getAgentConfigs(projectId)
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
  const sessionEventsQuery = useQuery({
    queryKey: ["runtime-events", selectedSessionId],
    queryFn: () => getRuntimeEvents(projectId, selectedSessionId!),
    enabled: Boolean(selectedSessionId)
  });
  const sessionTranscriptQuery = useQuery({
    queryKey: ["runtime-transcript", selectedSessionId],
    queryFn: () => getTranscript(projectId, selectedSessionId!),
    enabled: Boolean(selectedSessionId)
  });
  const sessionAuditsQuery = useQuery({
    queryKey: ["runtime-audits", selectedSessionId],
    queryFn: () => getTranscriptAccessAudits(projectId, selectedSessionId!),
    enabled: Boolean(selectedSessionId)
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
      setTaskActionMessage("任务已领取");
      await queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
    }
  });
  const taskMutation = useMutation({
    mutationFn: (payload: { title: string; type: string; assignedAgentId?: string }) =>
      createTask(projectId, missionId, payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
    }
  });

  const handoffMutation = useMutation({
    mutationFn: (taskId: string) => sendTaskToAdmin(projectId, missionId, taskId),
    onSuccess: async () => {
      setTaskActionMessage("任务已提交给管理员");
      await queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
    }
  });

  const checkpointMutation = useMutation({
    mutationFn: (taskId: string) => requestReviewCheckpoint(projectId, missionId, taskId),
    onSuccess: () => {
      setTaskActionMessage("已发起检查点");
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
  const selectedTask =
    taskBoardQuery.data?.items?.find((item) => item.id === selectedTaskId) ??
    taskBoardQuery.data?.items?.[0] ??
    null;
  const runtimeItems = runtimeQuery.data ?? [];

  useEffect(() => {
    return subscribeToChannel(projectId, `mission:${missionId}`, () => {
      void queryClient.invalidateQueries({ queryKey: ["mission-runtimes", missionId] });
      void queryClient.invalidateQueries({ queryKey: ["mission-task-board", missionId] });
      void queryClient.invalidateQueries({ queryKey: ["mission-approvals", missionId] });
      void queryClient.invalidateQueries({ queryKey: ["mission-archive", missionId] });
    });
  }, [missionId, queryClient]);

  useEffect(() => {
    if (!selectedTaskId && taskBoardQuery.data?.items?.length) {
      setSelectedTaskId(taskBoardQuery.data.items[0].id);
    }
  }, [selectedTaskId, taskBoardQuery.data]);

  useEffect(() => {
    if (!selectedSessionId && runtimeItems.length) {
      setSelectedSessionId(runtimeItems[0].currentExecutorSessionId ?? null);
    }
  }, [runtimeItems, selectedSessionId]);

  return (
    <PageShell
      eyebrow="Mission 工作台"
      title="任务作业空间"
      description="围绕当前 Mission 的任务推进、任务详情、Agent 执行日志和产物协同来组织页面。"
      aside={
        <div className="hero-stat">
          <span>工作区</span>
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
      <TaskComposer
        agents={agentsQuery.data ?? []}
        isSubmitting={taskMutation.isPending}
        onCreateTask={(payload) => taskMutation.mutate(payload)}
      />
      <TaskDetailPanel task={selectedTask} documents={documentsQuery.data ?? []} />
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
        onSelectTask={(taskId) => setSelectedTaskId(taskId)}
        onClaimTask={(taskId) => claimMutation.mutate(taskId)}
        onSendToAdmin={(taskId) => handoffMutation.mutate(taskId)}
        onRequestReview={(taskId) => checkpointMutation.mutate(taskId)}
      />
      <AgentRuntimeBoard
        runtimes={runtimeItems}
        onSelectRuntime={(sessionId) => setSelectedSessionId(sessionId)}
      />
      <TaskRuntimePanel
        runtimes={runtimeItems}
        events={sessionEventsQuery.data ?? []}
        transcript={sessionTranscriptQuery.data ?? null}
        audits={sessionAuditsQuery.data ?? []}
      />
    </PageShell>
  );
}
