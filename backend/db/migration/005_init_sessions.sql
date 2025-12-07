--TABLES
CREATE TABLE IF NOT EXISTS sessions (
  id uuid PRIMARY KEY,
  user_login TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  client_ip TEXT NOT NULL,
  -- 0 = false (not locked), 1 = true (locked)
  is_blocked INTEGER DEFAULT 0 NOT NULL,
  expires_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_login) REFERENCES users (user_login)
);
