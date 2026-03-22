import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { useEffect } from "react";

import { PageShell } from "../components/page-shell";
import { getRuntimeEvents, getRuntimeSession, getTranscript, getTranscriptAccessAudits } from "../lib/api";
import { useQueryClient } from "@tanstack/react-query";
import { subscribeToChannel } from "../lib/sse";
import { TranscriptAccessAuditTable } from "../features/runtime/transcript-access-audit-table";
import { RuntimeEventTimeline } from "../features/runtime/runtime-event-timeline";
import { TranscriptInspector } from "../features/runtime/transcript-inspector";

const projectId = "proj_1";

/**
 * RuntimeSessionPage 渲染当前路由对应的页面级工作区。
 */
export function RuntimeSessionPage() {
  const { sessionId = "session_1" } = useParams();
  const queryClient = useQueryClient();

  const sessionQuery = useQuery({
    queryKey: ["runtime-session", sessionId],
    queryFn: () => getRuntimeSession(projectId, sessionId)
  });
  const eventsQuery = useQuery({
    queryKey: ["runtime-events", sessionId],
    queryFn: () => getRuntimeEvents(projectId, sessionId)
  });
  const transcriptQuery = useQuery({
    queryKey: ["runtime-transcript", sessionId],
    queryFn: () => getTranscript(projectId, sessionId)
  });
  const auditQuery = useQuery({
    queryKey: ["runtime-audits", sessionId],
    queryFn: () => getTranscriptAccessAudits(projectId, sessionId)
  });

  useEffect(() => {
    return subscribeToChannel(projectId, `session:${sessionId}`, () => {
      void queryClient.invalidateQueries({ queryKey: ["runtime-session", sessionId] });
      void queryClient.invalidateQueries({ queryKey: ["runtime-events", sessionId] });
      void queryClient.invalidateQueries({ queryKey: ["runtime-transcript", sessionId] });
      void queryClient.invalidateQueries({ queryKey: ["runtime-audits", sessionId] });
    });
  }, [queryClient, sessionId]);

  return (
    <PageShell
      eyebrow="运行会话"
      title={sessionQuery.data?.backend ? `运行会话 · ${sessionQuery.data.backend}` : "运行会话详情"}
      description="优先展示结构化事件，必要时再查看脱敏后的 transcript 和访问审计。"
      aside={
        <div className="hero-stat">
          <span>状态</span>
          <strong>{sessionQuery.data?.status ?? "离线快照"}</strong>
        </div>
      }
    >
      <RuntimeEventTimeline events={eventsQuery.data ?? []} />
      <TranscriptInspector transcript={transcriptQuery.data ?? null} />
      <TranscriptAccessAuditTable audits={auditQuery.data ?? []} />
    </PageShell>
  );
}
