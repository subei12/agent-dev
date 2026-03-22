const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "";

/**
 * subscribeToChannel 为前端订阅实时 SSE 通道。
 */
export function subscribeToChannel(
  projectId: string,
  channel: string,
  onMessage: (event: MessageEvent<string>) => void
) {
  const path = `${API_BASE}/api/projects/${projectId}/events/stream?channel=${encodeURIComponent(channel)}`;
  const source = new EventSource(path);
  source.onmessage = onMessage;

  return () => {
    source.close();
  };
}
