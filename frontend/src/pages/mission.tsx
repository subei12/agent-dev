import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

import { PageShell } from "../components/page-shell";
import { getDiscussionSessions, getMission, getMissionDocuments, getMissionRuntimes, getTaskBoard } from "../lib/api";
import { DiscussionSessionPanel } from "../features/discussions/discussion-session-panel";
import { DocumentList } from "../features/documents/document-list";
import { MissionOverview } from "../features/missions/mission-overview";
import { TaskBoard } from "../features/tasks/task-board";
import { AgentRuntimeBoard } from "../features/runtime/agent-runtime-board";
import { RuntimeEventTimeline } from "../features/runtime/runtime-event-timeline";

const projectId = "proj_1";

export function MissionPage() {
  const { missionId = "mission_1" } = useParams();

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
      <DocumentList documents={documentsQuery.data ?? []} />
      <DiscussionSessionPanel sessions={sessionsQuery.data ?? []} />
      <TaskBoard tasks={taskBoardQuery.data?.items} />
      <AgentRuntimeBoard runtimes={runtimeQuery.data ?? []} />
      <RuntimeEventTimeline events={[]} />
    </PageShell>
  );
}
