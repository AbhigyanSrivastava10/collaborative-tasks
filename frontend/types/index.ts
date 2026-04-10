export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
  provider: string;
  created_at: string;
  updated_at: string;
}

export interface Board {
  id: string;
  name: string;
  description: string;
  owner_id: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  board_id: string;
  assigned_to: string | null;
  title: string;
  description: string;
  status: "todo" | "in_progress" | "done";
  priority: "low" | "medium" | "high";
  position: number;
  due_date: string | null;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
