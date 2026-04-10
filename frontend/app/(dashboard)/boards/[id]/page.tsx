"use client";

import { useQuery } from "@tanstack/react-query";
import { use } from "react";
import api from "@/lib/api";
import { Board, Task } from "@/types";
import KanbanBoard from "@/components/KanbanBoard";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { useWebSocket } from "@/hooks/useWebSocket";

export default function BoardPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);

  const { data: board, isLoading: boardLoading } = useQuery<Board>({
    queryKey: ["board", id],
    queryFn: async () => {
      const res = await api.get(`/api/boards/${id}`);
      return res.data;
    },
  });

  const { data: tasks = [], isLoading: tasksLoading } = useQuery<Task[]>({
    queryKey: ["tasks", id],
    queryFn: async () => {
      const res = await api.get(`/api/boards/${id}/tasks`);
      return res.data;
    },
  });
  useWebSocket(id);

  if (boardLoading || tasksLoading) {
    return <p className="text-slate-500">Loading board...</p>;
  }

  return (
    <div>
      <div className="mb-6">
        <Link href="/boards" className="flex items-center gap-1 text-sm text-slate-500 hover:text-slate-700 mb-3">
          <ArrowLeft className="w-4 h-4" />
          Back to boards
        </Link>
        <h1 className="text-2xl font-bold text-slate-900">{board?.name}</h1>
        {board?.description && (
          <p className="text-slate-500 mt-1">{board.description}</p>
        )}
      </div>
      <KanbanBoard boardID={id} initialTasks={tasks} />
    </div>
  );
}
