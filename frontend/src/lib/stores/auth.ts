import { user } from "./user";

let accessToken: string | null = null;

export async function login(user_login: string, password: string) {
  const res = await fetch("/users/login", {
    method: "POST",
    body: JSON.stringify({ user_login, password }),
    headers: { "Content-Type": "application/json" },
  });

  if (!res.ok) throw new Error("Login failed");

  const data = await res.json();
  accessToken = data.access_token;
  user.set(data.user);
}

export async function logout() {
  await fetch("/api/auth/logout", { method: "POST" });
  user.set(null);
  accessToken = null;
}

export async function authFetch(url: string, options: RequestInit = {}) {
  // attach access token
  options.headers = {
    ...(options.headers || {}),
    Authorization: `Bearer ${accessToken}`,
  };

  const res = await fetch(url, options);

  // If expired, refresh token
  if (res.status === 401) {
    const refreshed = await refreshToken();
    if (!refreshed) {
      logout();
      throw new Error("Session expired");
    }

    // retry request with new token
    return authFetch(url, options);
  }

  return res;
}

async function refreshToken() {
  const res = await fetch("/api/auth/refresh", {
    method: "POST",
    credentials: "include", // includes HTTP-only cookie
  });

  if (!res.ok) return false;

  const data = await res.json();
  accessToken = data.access_token;
  user.set(data.user);
  return true;
}
