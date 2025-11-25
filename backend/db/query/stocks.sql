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
  stock_data (
    stock_id,
    event_label,
    value_dollars,
    value_cents,
    sequence
  )
VALUES
  (?, ?, ?, ?, ?);

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
      sd.sequence DESC
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
  sequence ASC
LIMIT
  ?;
