import { writable } from "svelte/store";

export interface User {
  id: number;
  user_login: string;
  role: "user" | "admin";
  dollars: number;
  cents: number;
}

export interface SessionInfo {
  session_id: string;
  access_token: string;
  access_token_expires_at: string; // ISO timestamp from Go
  refresh_token: string;
  refresh_token_expires_at: string;
}

export const user = writable<User | null>(null);
export const session = writable<SessionInfo | null>(null);
