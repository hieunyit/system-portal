import { useEffect, useRef } from 'react';

export function useWebSocket(url: string) {
  const socket = useRef<WebSocket | null>(null);
  useEffect(() => {
    socket.current = new WebSocket(url);
    return () => socket.current?.close();
  }, [url]);
  return socket.current;
}
