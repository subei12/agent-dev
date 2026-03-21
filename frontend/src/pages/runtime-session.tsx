import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";

import { PageShell } from "../components/page-shell";
import { getRuntimeEvents, getRuntimeSession, getTranscript, getTranscriptAccessAudits } from "../lib/api";
import { TranscriptAccessAuditTable } from "../features/runtime/transcript-access-audit-table";
import { RuntimeEventTimeline } from "../features/runtime/runtime-event-timeline";
import { TranscriptInspector } from "../features/runtime/transcript-inspector";

const projectId = "proj_1";

export function RuntimeSessionPage() {
  const { sessionId = "session_1" } = useParams();

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

  return (
    <PageShell
      eyebrow="Session Inspector"
      title={sessionQuery.data?.backend ? `Runtime Session · ${sessionQuery.data.backend}` : "Runtime Session Inspector"}
      description="Structured events stay primary. Redacted transcript expands underneath when someone needs the actual model dialogue."
      aside={
        <div className="hero-stat">
          <span>Status</span>
          <strong>{sessionQuery.data?.status ?? "offline snapshot"}</strong>
        </div>
      }
    >
      <RuntimeEventTimeline events={eventsQuery.data ?? []} />
      <TranscriptInspector transcript={transcriptQuery.data ?? null} />
      <TranscriptAccessAuditTable audits={auditQuery.data ?? []} />
    </PageShell>
  );
}
