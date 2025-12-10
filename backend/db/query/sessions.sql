-- name: CreateSession :one
INSERT INTO
  sessions (
    id,
    user_login,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: GetSession :one
SELECT
  *
FROM
  sessions
WHERE
  user_login = ?
LIMIT
  1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
  user_login = ?;
