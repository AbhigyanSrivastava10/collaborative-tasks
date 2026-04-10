"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { Board } from "@/types";
import BoardCard from "@/components/BoardCard";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Plus } from "lucide-react";

export default function BoardsPage() {
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [open, setOpen] = useState(false);

  const { data: boards = [], isLoading } = useQuery<Board[]>({
    queryKey: ["boards"],
    queryFn: async () => {
      const res = await api.get("/api/boards");
      return res.data;
    },
  });

  const createBoard = useMutation({
    mutationFn: async () => {
      const res = await api.post("/api/boards", { name, description });
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["boards"] });
      setName("");
      setDescription("");
      setOpen(false);
    },
  });

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-slate-900">My Boards</h1>
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              New Board
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create a new board</DialogTitle>
            </DialogHeader>
            <div className="space-y-4 pt-2">
              <div className="space-y-2">
                <Label>Name</Label>
                <Input
                  placeholder="e.g. Product Roadmap"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label>Description</Label>
                <Input
                  placeholder="Optional description"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>
              <Button
                className="w-full"
                onClick={() => createBoard.mutate()}
                disabled={!name || createBoard.isPending}
              >
                {createBoard.isPending ? "Creating..." : "Create Board"}
              </Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      {isLoading ? (
        <p className="text-slate-500">Loading boards...</p>
      ) : boards.length === 0 ? (
        <div className="text-center py-20 text-slate-400">
          <p className="text-lg font-medium">No boards yet</p>
          <p className="text-sm mt-1">Create your first board to get started</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {boards.map((board) => (
            <BoardCard key={board.id} board={board} />
          ))}
        </div>
      )}
    </div>
  );
}
