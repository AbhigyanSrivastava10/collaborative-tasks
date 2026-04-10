"use client";

import { useState } from "react";
import { DragDropContext, Droppable, Draggable, DropResult } from "@hello-pangea/dnd";
import { Task } from "@/types";
import TaskCard from "./TaskCard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";
import api from "@/lib/api";

const COLUMNS = [
  { id: "todo", label: "To Do" },
  { id: "in_progress", label: "In Progress" },
  { id: "done", label: "Done" },
];

export default function KanbanBoard({ boardID, initialTasks }: { boardID: string; initialTasks: Task[] }) {
  const [tasks, setTasks] = useState<Task[]>(initialTasks);
  const [newTaskTitle, setNewTaskTitle] = useState<Record<string, string>>({});
  const [adding, setAdding] = useState<Record<string, boolean>>({});

  function getTasksByStatus(status: string) {
    return tasks.filter((t) => t.status === status);
  }

  async function handleDragEnd(result: DropResult) {
    const { destination, source, draggableId } = result;
    if (!destination) return;
    if (destination.droppableId === source.droppableId && destination.index === source.index) return;

    const task = tasks.find((t) => t.id === draggableId);
    if (!task) return;

    // Optimistically update UI
    const updated = tasks.map((t) =>
      t.id === draggableId
        ? { ...t, status: destination.droppableId as Task["status"], position: destination.index }
        : t
    );
    setTasks(updated);

    // Persist to backend
    try {
      await api.put(`/api/boards/${boardID}/tasks/${draggableId}`, {
        ...task,
        status: destination.droppableId,
        position: destination.index,
      });
    } catch {
      setTasks(tasks); // revert on error
    }
  }

  async function handleAddTask(status: string) {
    const title = newTaskTitle[status]?.trim();
    if (!title) return;

    try {
      const res = await api.post(`/api/boards/${boardID}/tasks`, { title, status });
      setTasks((prev) => [...prev, res.data]);
      setNewTaskTitle((prev) => ({ ...prev, [status]: "" }));
      setAdding((prev) => ({ ...prev, [status]: false }));
    } catch (err) {
      console.error("Failed to create task", err);
    }
  }

  return (
    <DragDropContext onDragEnd={handleDragEnd}>
      <div className="flex gap-4 overflow-x-auto pb-4">
        {COLUMNS.map((col) => (
          <div key={col.id} className="flex-shrink-0 w-72">
            <div className="bg-slate-100 rounded-lg p-3">
              <div className="flex items-center justify-between mb-3">
                <h3 className="font-semibold text-slate-700 text-sm">{col.label}</h3>
                <span className="text-xs text-slate-400 bg-slate-200 px-2 py-0.5 rounded-full">
                  {getTasksByStatus(col.id).length}
                </span>
              </div>

              <Droppable droppableId={col.id}>
                {(provided) => (
                  <div
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                    className="min-h-[100px]"
                  >
                    {getTasksByStatus(col.id).map((task, index) => (
                      <Draggable key={task.id} draggableId={task.id} index={index}>
                        {(provided) => (
                          <div
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            {...provided.dragHandleProps}
                          >
                            <TaskCard task={task} />
                          </div>
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>

              {adding[col.id] ? (
                <div className="mt-2 space-y-2">
                  <Input
                    autoFocus
                    placeholder="Task title..."
                    value={newTaskTitle[col.id] || ""}
                    onChange={(e) => setNewTaskTitle((prev) => ({ ...prev, [col.id]: e.target.value }))}
                    onKeyDown={(e) => e.key === "Enter" && handleAddTask(col.id)}
                  />
                  <div className="flex gap-2">
                    <Button size="sm" onClick={() => handleAddTask(col.id)}>Add</Button>
                    <Button size="sm" variant="ghost" onClick={() => setAdding((prev) => ({ ...prev, [col.id]: false }))}>
                      Cancel
                    </Button>
                  </div>
                </div>
              ) : (
                <Button
                  variant="ghost"
                  size="sm"
                  className="w-full mt-2 text-slate-500 hover:text-slate-700"
                  onClick={() => setAdding((prev) => ({ ...prev, [col.id]: true }))}
                >
                  <Plus className="w-4 h-4 mr-1" /> Add task
                </Button>
              )}
            </div>
          </div>
        ))}
      </div>
    </DragDropContext>
  );
}

