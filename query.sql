-- name: GetOLHC :one
SELECT * FROM tbl_ohlc
WHERE day = $1 and ticker_id = $2 LIMIT 1;

-- name: ListOLHCs :many
SELECT * FROM tbl_ohlc
ORDER BY day;

-- name: CreateOLHC :one
INSERT INTO tbl_ohlc (
  day, ticker_id, o, h, l, c, v
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateOLHC :exec
UPDATE tbl_ohlc
  set o = $3,
  h = $4,
  l = $5,
  c = $6,
  v = $7
WHERE day = $1 and ticker_id = $2;

-- name: DeleteOLHC :exec
DELETE FROM tbl_ohlc
WHERE day = $1 and ticker_id = $2;
