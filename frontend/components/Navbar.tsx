"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import { LayoutDashboard, LogOut } from "lucide-react";

export default function Navbar() {
  const router = useRouter();
  const { user, clearAuth } = useAuthStore();

  function handleLogout() {
    clearAuth();
    router.push("/login");
  }

  return (
    <nav className="bg-white border-b border-slate-200 px-4 py-3">
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <Link href="/boards" className="flex items-center gap-2 font-semibold text-slate-900">
          <LayoutDashboard className="w-5 h-5" />
          Collaborative Tasks
        </Link>
        <div className="flex items-center gap-4">
          <span className="text-sm text-slate-500">{user?.name}</span>
          <Button variant="ghost" size="sm" onClick={handleLogout}>
            <LogOut className="w-4 h-4 mr-2" />
            Logout
          </Button>
        </div>
      </div>
    </nav>
  );
}

