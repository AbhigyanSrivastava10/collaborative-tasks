import api from "./api";
import { AuthResponse } from "@/types";

export async function login(email: string, password: string): Promise<AuthResponse> {
  const res = await api.post<AuthResponse>("/api/auth/login", { email, password });
  return res.data;
}

export async function register(email: string, password: string, name: string): Promise<AuthResponse> {
  const res = await api.post<AuthResponse>("/api/auth/register", { email, password, name });
  return res.data;
}
