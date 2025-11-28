-- name: CreateUserStock :one
INSERT INTO
  user_stocks (user_id, stock_id, quantity)
VALUES
  (?, ?, ?) RETURNING *;

-- name: GetUserStock :one
SELECT
  *
FROM
  user_stocks
WHERE
  user_id = ?
  AND stock_id = ?
LIMIT
  1;

-- name: UpdateUserStock :one
UPDATE user_stocks
SET
  quantity = sqlc.arg (quantity)
WHERE
  user_id = sqlc.arg (user_id)
  AND stock_id = sqlc.arg (stock_id) RETURNING *;

-- name: DeleteUserStock :exec
DELETE FROM user_stocks
WHERE
  user_id = ?
  AND stock_id = ?
  AND quantity = 0;
