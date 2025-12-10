import { user, session, type SessionInfo, type User } from "./user";

let _session: SessionInfo | null = null;
session.subscribe((v) => (_session = v));

let _user: User | null = null;
user.subscribe((v) => {
  _user = v;
});

export async function login(user_login: string, password: string) {
  const res = await fetch("users/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_login, password }),
  });

  if (!res.ok) {
    throw new Error("Invalid username or password");
  }

  const data = await res.json();

  const newSession: SessionInfo = {
    session_id: data.session_id,
    access_token: data.access_token,
    access_token_expires_at: data.access_token_expires_at,
    refresh_token: data.refresh_token,
    refresh_token_expires_at: data.refresh_token_expires_at,
  };

  session.set(newSession);
  user.set(data.user);
}

export async function logout(user_login: string) {
  if (_session) {
    await fetch("users/logout", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${_session.access_token}`,
      },
      body: JSON.stringify({ user_login }),
    });
  }

  session.set(null);
  user.set(null);
}

export async function authFetch(url: string, options: RequestInit = {}) {
  if (!_session) {
    throw new Error("No active session");
  }
  if (!_user) {
    throw new Error("no user information");
  }

  // attach access token
  options.headers = {
    ...(options.headers || {}),
    Authorization: `Bearer ${_session.access_token}`,
  };

  const res = await fetch(url, options);

  // Access token expired → try refresh
  if (res.status === 401) {
    const didRefresh = await attemptRefresh();

    if (!didRefresh) {
      await logout(_user.user_login);
      throw new Error("Session expired");
    }

    // retry the request
    if (!_session) throw new Error("No session after refresh");

    options.headers = {
      ...(options.headers || {}),
      Authorization: `Bearer ${_session.access_token}`,
    };

    return fetch(url, options);
  }

  return res;
}

async function attemptRefresh(): Promise<boolean> {
  if (!_session) return false;

  const res = await fetch("/api/auth/refresh", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${_session.refresh_token}`,
    },
  });

  if (!res.ok) return false;

  const data = await res.json();

  const newSession: SessionInfo = {
    session_id: _session.session_id,
    access_token: data.access_token,
    access_token_expires_at: data.access_token_expires_at,
    refresh_token: _session.refresh_token,
    refresh_token_expires_at: _session.refresh_token_expires_at,
  };

  session.set(newSession);

  return true;
}
