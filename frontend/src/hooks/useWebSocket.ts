import { useEffect, useRef, useCallback } from 'react';
import type { WebSocketMessage, GlobeUpdate } from '../types';
import { useGlobeStore } from '../store';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080';

export const useWebSocket = () => {
  const wsRef = useRef<WebSocket | null>(null);
  const { updateCityScore } = useGlobeStore();

  const connect = useCallback(() => {
    const ws = new WebSocket(`${WS_URL}/ws/globe`);

    ws.onopen = () => {
      console.log('WebSocket connected');
    };

    ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);

        switch (message.type) {
          case 'connected':
            console.log('WebSocket connection established');
            break;

          case 'city_update':
            const update = message.data as GlobeUpdate;
            updateCityScore(update.city_id, update.new_score);
            console.log(`City ${update.city_name} score updated to ${update.new_score}`);
            break;

          default:
            console.log('Unknown message type:', message.type);
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected, reconnecting in 5s...');
      setTimeout(connect, 5000);
    };

    wsRef.current = ws;
  }, [updateCityScore]);

  useEffect(() => {
    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  return {
    isConnected: wsRef.current?.readyState === WebSocket.OPEN,
  };
};
