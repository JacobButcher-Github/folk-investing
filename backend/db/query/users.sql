-- name: CreateUser :one
INSERT INTO
  users (user_login, hashed_password, dollars, cents)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: GetUserFromId :one
SELECT
  *
FROM
  users
WHERE
  id = ?
LIMIT
  1;

-- name: GetUserFromName :one
SELECT
  *
FROM
  users
WHERE
  user_login = ?
LIMIT
  1;

-- name: UpdateUser :one
UPDATE users
SET
  user_login = COALESCE(sqlc.narg (new_login), user_login),
  hashed_password = COALESCE(sqlc.narg (hashed_password), hashed_password),
  dollars = COALESCE(sqlc.narg (dollars), dollars),
  cents = COALESCE(sqlc.narg (cents), cents)
WHERE
  user_login = sqlc.arg (user_login) RETURNING *;
