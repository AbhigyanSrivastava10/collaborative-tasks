import Link from "next/link";
import { Board } from "@/types";
import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { LayoutDashboard } from "lucide-react";

export default function BoardCard({ board }: { board: Board }) {
  return (
    <Link href={`/boards/${board.id}`}>
      <Card className="hover:shadow-md transition-shadow cursor-pointer h-full">
        <CardHeader>
          <div className="flex items-center gap-2 mb-1">
            <LayoutDashboard className="w-4 h-4 text-slate-400" />
            <CardTitle className="text-lg">{board.name}</CardTitle>
          </div>
          <CardDescription>{board.description || "No description"}</CardDescription>
        </CardHeader>
      </Card>
    </Link>
  );
}

