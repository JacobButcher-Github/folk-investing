-- name: CreateStock :one
INSERT INTO
  stocks (name, image_path)
VALUES
  (?, ?) RETURNING *;

-- name: GetStockFromId :one
SELECT
  *
FROM
  stocks
WHERE
  id = ?
LIMIT
  1;

-- name: GetStockFromName :one
SELECT
  *
FROM
  stocks
WHERE
  name = ?
LIMIT
  1;

-- name: GetAllStocks :many
SELECT
  *
FROM
  STOCKS;

-- name: UpdateStock :one
UPDATE stocks
SET
  name = COALESCE(sqlc.narg (new_name), name),
  image_path = COALESCE(sqlc.narg (image_path), image_path)
WHERE
  name = sqlc.arg (name) RETURNING *;

-- name: DeleteStock :exec
DELETE FROM stocks
WHERE
  name = ?;

-- name: CreateStockData :one
INSERT INTO
  stock_data (stock_id, event_label, value_dollars, value_cents)
VALUES
  (?, ?, ?, ?) RETURNING *;

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

-- name: GetStocksData :many
-- Get all data for stocks up to a limit. Meant for use by the graph  on the main page.
SELECT
  sd.*
FROM
  stock_data sd
WHERE
  sd.stock_id IN (sqlc.slice (stock_ids))
  AND sd.id IN (
    SELECT
      id
    FROM
      stock_data
    WHERE
      stock_id = sd.stock_id
    ORDER BY
      id DESC
    LIMIT
      ?
  )
ORDER BY
  sd.stock_id,
  sd.id ASC;

-- name: GetStockDataByLabel :many
-- Get all data for stocks with a certain label.
SELECT
  *
FROM
  stock_data
WHERE
  event_label = ?;

-- name: DeleteStockDataByLabel :exec
-- Delete all stock data associated with a certain label.
DELETE FROM stock_data
WHERE
  event_label = ?;

-- name: UpdateStockData :one
UPDATE stock_data
SET
  stock_id = COALESCE(sqlc.narg (new_id), stock_id),
  event_label = COALESCE(sqlc.narg (new_label), event_label),
  value_dollars = COALESCE(sqlc.narg (value_dollars), value_dollars),
  value_cents = COALESCE(sqlc.narg (value_cents), value_cents)
WHERE
  stock_id = sqlc.arg (stock_id)
  AND event_label = sqlc.arg (event_label) RETURNING *;
