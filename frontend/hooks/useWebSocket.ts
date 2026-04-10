"use client";

import { useEffect, useRef } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Task } from "@/types";

interface WSMessage {
  board_id: string;
  type: "task_created" | "task_updated" | "task_deleted";
  payload: Task;
}

export function useWebSocket(boardID: string) {
  const queryClient = useQueryClient();
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!boardID) return;

    const token = localStorage.getItem("token");
    const wsURL = `ws://localhost:8080/ws/boards/${boardID}?token=${token}`;
    const ws = new WebSocket(wsURL);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("ws: connected to board", boardID);
    };

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);

        queryClient.setQueryData<Task[]>(["tasks", boardID], (old = []) => {
          switch (msg.type) {
            case "task_created":
              // Add task if it doesn't already exist
              if (old.find((t) => t.id === msg.payload.id)) return old;
              return [...old, msg.payload];

            case "task_updated":
              return old.map((t) => (t.id === msg.payload.id ? msg.payload : t));

            case "task_deleted":
              return old.filter((t) => t.id !== msg.payload.id);

            default:
              return old;
          }
        });
      } catch (err) {
        console.error("ws: failed to parse message", err);
      }
    };

    ws.onerror = (err) => {
      console.error("ws: error", err);
    };

    ws.onclose = () => {
      console.log("ws: disconnected from board", boardID);
    };

    // Cleanup on unmount
    return () => {
      ws.close();
    };
  }, [boardID, queryClient]);

  return wsRef;
}
