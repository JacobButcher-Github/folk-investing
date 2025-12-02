-- name: CreateSiteSettings :one
INSERT INTO
  site_settings (
    number_of_events_visible,
    value_symbol,
    event_label,
    lockout_time_start
  )
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: GetSiteSettings :one
SELECT
  *
FROM
  site_settings
LIMIT
  1;

-- name: GetNumberEvents :one
SELECT
  number_of_events_visible
FROM
  site_settings
LIMIT
  1;

-- name: GetValueSymbol :one
SELECT
  value_symbol
FROM
  site_settings
LIMIT
  1;

-- name: GetEventLabel :one
SELECT
  event_label
FROM
  site_settings
LIMIT
  1;

-- name: GetLockoutStatus :one
SELECT
  lockout
FROM
  site_settings
LIMIT
  1;

-- name: GetLockoutTime :one
SELECT
  lockout_time_start
FROM
  site_settings
LIMIT
  1;

-- name: UpdateSettings :one
UPDATE site_settings
SET
  number_of_events_visible = COALESCE(
    sqlc.narg (number_of_events_visible),
    number_of_events_visible
  ),
  value_symbol = COALESCE(sqlc.narg (value_symbol), value_symbol),
  event_label = COALESCE(sqlc.narg (event_label), event_label),
  lockout_time_start = COALESCE(
    sqlc.narg (lockout_time_start),
    lockout_time_start
  )
WHERE
  id = 0 RETURNING *;
