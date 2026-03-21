export function subscribe(path: string, onMessage: (event: MessageEvent<string>) => void) {
  const source = new EventSource(path);
  source.onmessage = onMessage;

  return () => {
    source.close();
  };
}
