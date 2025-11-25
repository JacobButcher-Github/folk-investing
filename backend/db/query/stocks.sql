-- name: CreateStock :exec
INSERT INTO
  stocks (name, image_path)
VALUES
  (?, ?);

-- name: GetStock :one
SELECT
  *
FROM
  stocks
WHERE
  name = ?
LIMIT
  1;

-- name: UpdateStock :one
UPDATE stocks
SET
  name = COALESCE(sqlc.narg (new_name), name),
  image_path = COALESCE(sqlc.narg (image_path), image_path)
WHERE
  name = sqlc.arg (name) RETURNING *;

-- name: CreateStockData :exec
INSERT INTO
  stock_data (stock_id, event_label, value_dollars, value_cents)
VALUES
  (?, ?, ?, ?);

-- name: PruneStockData :exec
DELETE FROM stock_data
WHERE
  stock_data.stock_id = ?
  AND stock_data.id NOT IN (
    SELECT
      sd.id
    FROM
      stock_data as sd
    WHERE
      sd.stock_id = ?
    ORDER BY
      sd.id DESC
    LIMIT
      ?
  );

-- name: GetStockData :one
SELECT
  *
FROM
  stock_data
WHERE
  stock_id = ?
ORDER BY
  id ASC
LIMIT
  ?;

-- name: UpdateStockData :one
UPDATE stock_data
SET
  stock_id = COALESCE(sqlc.narg (new_id), stock_id),
  stock_id = COALESCE(sqlc.narg (new_label), event_label),
  stock_id = COALESCE(sqlc.narg (value_dollars), value_dollars),
  stock_id = COALESCE(sqlc.narg (value_cents), value_cents)
WHERE
  stock_id = sqlc.arg (stock_id)
  AND event_label = sqlc.arg (event_label) RETURNING *;
