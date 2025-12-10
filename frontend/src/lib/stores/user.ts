import { writable } from "svelte/store";

export interface User {
  id: number;
  user_login: string;
  role: "user" | "admin";
}

export const user = writable<User | null>(null);
