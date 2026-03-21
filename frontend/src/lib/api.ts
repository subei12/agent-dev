export type Mission = {
  id: string;
  title: string;
  status: string;
  projectId?: string;
  description?: string;
};

export type DocumentItem = {
  id: string;
  title: string;
  kind: string;
  currentAdoptedVersionId?: string;
};

export type DiscussionSession = {
  id: string;
  topic: string;
  status: string;
};

export type TaskBoard = {
  id: string;
  title: string;
  items?: Array<{
    id: string;
    title: string;
    type: string;
    status: string;
    assignedAgentId?: string;
  }>;
};

export type MissionRuntime = {
  id: string;
  agentId: string;
  status: string;
  statusSummary?: string;
  currentExecutorSessionId?: string;
};

export type RuntimeEvent = {
  id: string;
  title: string;
  type: string;
  category: string;
  summary?: string;
};

export type TranscriptView = {
  transcript: {
    id: string;
    status: string;
    executorSessionId: string;
  };
  entries?: Array<{
    id: string;
    role: string;
    entryType: string;
    redactedText?: string;
  }>;
};

export type TranscriptAccessAudit = {
  id: string;
  actorUserId: string;
  accessMode: string;
  reason?: string;
};

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://127.0.0.1:8080";

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`);
  if (!response.ok) {
    throw new Error(`request failed: ${response.status}`);
  }
  return response.json() as Promise<T>;
}

export async function getMission(projectId: string, missionId: string): Promise<Mission | null> {
  try {
    return await request<Mission>(`/api/projects/${projectId}/missions/${missionId}`);
  } catch {
    return null;
  }
}

export async function getMissionDocuments(projectId: string, missionId: string): Promise<DocumentItem[]> {
  try {
    return await request<DocumentItem[]>(`/api/projects/${projectId}/missions/${missionId}/documents`);
  } catch {
    return [];
  }
}

export async function getDiscussionSessions(projectId: string, missionId: string): Promise<DiscussionSession[]> {
  try {
    return await request<DiscussionSession[]>(`/api/projects/${projectId}/missions/${missionId}/discussions`);
  } catch {
    return [];
  }
}

export async function getMissionRuntimes(projectId: string, missionId: string): Promise<MissionRuntime[]> {
  try {
    return await request<MissionRuntime[]>(`/api/projects/${projectId}/missions/${missionId}/agent-runtimes`);
  } catch {
    return [];
  }
}

export async function getTaskBoard(projectId: string, missionId: string): Promise<TaskBoard | null> {
  try {
    return await request<TaskBoard>(`/api/projects/${projectId}/missions/${missionId}/task-board`);
  } catch {
    return null;
  }
}

export async function getRuntimeSession(projectId: string, sessionId: string) {
  try {
    return await request<{ id: string; backend: string; status: string }>(
      `/api/projects/${projectId}/executor-sessions/${sessionId}`
    );
  } catch {
    return null;
  }
}

export async function getRuntimeEvents(projectId: string, sessionId: string): Promise<RuntimeEvent[]> {
  try {
    return await request<RuntimeEvent[]>(`/api/projects/${projectId}/executor-sessions/${sessionId}/events`);
  } catch {
    return [];
  }
}

export async function getTranscript(projectId: string, sessionId: string): Promise<TranscriptView | null> {
  try {
    return await request<TranscriptView>(
      `/api/projects/${projectId}/executor-sessions/${sessionId}/transcript?view=redacted`
    );
  } catch {
    return null;
  }
}

export async function getTranscriptAccessAudits(
  projectId: string,
  sessionId: string
): Promise<TranscriptAccessAudit[]> {
  try {
    return await request<TranscriptAccessAudit[]>(
      `/api/projects/${projectId}/executor-sessions/${sessionId}/access-audits`
    );
  } catch {
    return [];
  }
}
