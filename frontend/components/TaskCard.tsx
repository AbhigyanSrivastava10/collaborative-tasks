import { Task } from "@/types";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Calendar } from "lucide-react";

const priorityColors: Record<string, string> = {
  low: "bg-slate-100 text-slate-600",
  medium: "bg-yellow-100 text-yellow-700",
  high: "bg-red-100 text-red-600",
};

export default function TaskCard({ task }: { task: Task }) {
  return (
    <Card className="mb-2 cursor-grab active:cursor-grabbing hover:shadow-sm transition-shadow">
      <CardContent className="p-3 space-y-2">
        <p className="text-sm font-medium text-slate-900">{task.title}</p>
        {task.description && (
          <p className="text-xs text-slate-500 line-clamp-2">{task.description}</p>
        )}
        <div className="flex items-center justify-between">
          <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${priorityColors[task.priority]}`}>
            {task.priority}
          </span>
          {task.due_date && (
            <span className="text-xs text-slate-400 flex items-center gap-1">
              <Calendar className="w-3 h-3" />
              {new Date(task.due_date).toLocaleDateString()}
            </span>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

