const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://127.0.0.1:8080";

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
