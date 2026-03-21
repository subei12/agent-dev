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

export type DocumentVersion = {
  id: string;
  documentId: string;
  version: number;
  status: string;
  contentFormat: string;
  contentHash: string;
  contentText: string;
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

export type MissionArchiveResult = {
  archive: {
    id: string;
    missionId: string;
    status: string;
  };
  manifest: {
    missionId: string;
    objectKeys: string[];
  };
};

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://127.0.0.1:8080";
const ACTOR_ID = "frontend-demo-user";

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      "X-Actor-Id": ACTOR_ID
    }
  });
  if (!response.ok) {
    throw new Error(`request failed: ${response.status}`);
  }
  return response.json() as Promise<T>;
}

async function send<T>(path: string, method: string, body?: unknown): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      "X-Actor-Id": ACTOR_ID
    },
    body: body ? JSON.stringify(body) : undefined
  });

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
    const response = await fetch(
      `${API_BASE}/api/projects/${projectId}/executor-sessions/${sessionId}/transcript?view=redacted`,
      {
        headers: {
          "X-Actor-Id": ACTOR_ID
        }
      }
    );
    if (!response.ok) {
      throw new Error(`request failed: ${response.status}`);
    }
    return response.json() as Promise<TranscriptView>;
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

/**
 * createDiscussionSession 向后端 API 发送写操作。
 */
export function createDiscussionSession(projectId: string, missionId: string, topic: string) {
  return send<DiscussionSession>(`/api/projects/${projectId}/missions/${missionId}/discussions`, "POST", {
    topic,
    initiatedByAgentId: "agent_admin"
  });
}

/**
 * createDocument 向后端 API 发送写操作。
 */
export function createDocument(projectId: string, missionId: string, title: string, kind: string) {
  return send<DocumentItem>(`/api/projects/${projectId}/missions/${missionId}/documents`, "POST", {
    title,
    kind
  });
}

export async function getDocumentVersions(
  projectId: string,
  missionId: string,
  documentId: string
): Promise<DocumentVersion[]> {
  try {
    return await request<DocumentVersion[]>(
      `/api/projects/${projectId}/missions/${missionId}/documents/${documentId}/versions`
    );
  } catch {
    return [];
  }
}

/**
 * createDocumentVersion 向后端 API 发送写操作。
 */
export function createDocumentVersion(
  projectId: string,
  missionId: string,
  documentId: string,
  contentText: string
) {
  const contentHash = `hash-${documentId}-${contentText.length}`;

  return send<DocumentVersion>(
    `/api/projects/${projectId}/missions/${missionId}/documents/${documentId}/versions`,
    "POST",
    {
      contentFormat: "md",
      storageKind: "db_text",
      contentHash,
      contentText,
      producedByAgentId: "agent_admin"
    }
  );
}

/**
 * adoptDocumentVersion 向后端 API 发送写操作。
 */
export function adoptDocumentVersion(
  projectId: string,
  missionId: string,
  documentId: string,
  versionId: string
) {
  return send<DocumentVersion>(
    `/api/projects/${projectId}/missions/${missionId}/documents/${documentId}/adopt`,
    "POST",
    {
      versionId
    }
  );
}

/**
 * claimTask 向后端 API 发送写操作。
 */
export function claimTask(projectId: string, missionId: string, taskId: string) {
  return send<{ id: string; status: string }>(
    `/api/projects/${projectId}/missions/${missionId}/tasks/${taskId}/claim`,
    "POST",
    {
      agentId: "agent_backend",
      claimReason: "workspace claim"
    }
  );
}

/**
 * sendTaskToAdmin 向后端 API 发送写操作。
 */
export function sendTaskToAdmin(projectId: string, missionId: string, taskId: string) {
  return send<{ id: string; status: string }>(
    `/api/projects/${projectId}/missions/${missionId}/tasks/${taskId}/handoffs`,
    "POST",
    {
      fromAgentId: "agent_backend",
      toAdminAgent: true,
      summary: "Task completed and awaiting admin review."
    }
  );
}

/**
 * requestReviewCheckpoint 向后端 API 发送写操作。
 */
export function requestReviewCheckpoint(projectId: string, missionId: string, taskId: string) {
  return send<{ id: string }>(
    `/api/projects/${projectId}/missions/${missionId}/tasks/${taskId}/review-checkpoints`,
    "POST",
    {
      kind: "doc_review",
      requestedByAgentId: "agent_admin",
      assignedAgentId: "agent_admin",
      summary: "Checkpoint requested from mission workspace."
    }
  );
}

/**
 * createArchive 向后端 API 发送写操作。
 */
export function createArchive(projectId: string, missionId: string) {
  return send<MissionArchiveResult>(`/api/projects/${projectId}/missions/${missionId}/archive`, "POST");
}
