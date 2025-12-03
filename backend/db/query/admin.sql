-- name: CreateAdmin :one
INSERT INTO
  users (user_login, role, hashed_password, dollars, cents)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

-- name: AdminUpdateUser :one
UPDATE users
SET
  role = COALESCE(sqlc.narg (role), role),
  dollars = COALESCE(sqlc.narg (dollars), dollars),
  cents = COALESCE(sqlc.narg (cents), cents)
WHERE
  user_login = sqlc.arg (user_login) RETURNING *
