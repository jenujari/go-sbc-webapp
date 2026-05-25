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

-- name: ListTickers :many
SELECT id, name, "desc", full_name, exchange
FROM tbl_ticker
ORDER BY name;

-- name: UpsertOLHC :execrows
INSERT INTO tbl_ohlc (
  day, ticker_id, o, h, l, c, v
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (day, ticker_id) DO UPDATE SET
  o = EXCLUDED.o,
  h = EXCLUDED.h,
  l = EXCLUDED.l,
  c = EXCLUDED.c,
  v = EXCLUDED.v;

-- name: UpsertPanchang :exec
INSERT INTO tbl_panchang (
  "time", tithy, nakshatra, weekday
) VALUES (
  $1, $2, $3, $4
)
ON CONFLICT ("time") DO UPDATE SET
  tithy = EXCLUDED.tithy,
  nakshatra = EXCLUDED.nakshatra,
  weekday = EXCLUDED.weekday;

-- name: UpsertPlanetPosition :exec
INSERT INTO tbl_planet_positions (
  observation_time, planet_name, longitude, latitude, distance,
  speed_long, speed_lat, speed_dist, speed_category, vedha,
  sign, nakshatra_name, nakshatra_pada, is_retro, sign_lord,
  sign_lordship, navamsa_sign, vargottama, uday_bala, uchcha_bala,
  vakra_bala, kshetra_bala, navamsha_bala, longitude_dms, latitude_dms,
  speed_long_dms
) VALUES (
  $1, $2, $3, $4, $5,
  $6, $7, $8, $9, $10,
  $11, $12, $13, $14, $15,
  $16, $17, $18, $19, $20,
  $21, $22, $23, $24, $25,
  $26
)
ON CONFLICT (observation_time, planet_name) DO UPDATE SET
  longitude = EXCLUDED.longitude,
  latitude = EXCLUDED.latitude,
  distance = EXCLUDED.distance,
  speed_long = EXCLUDED.speed_long,
  speed_lat = EXCLUDED.speed_lat,
  speed_dist = EXCLUDED.speed_dist,
  speed_category = EXCLUDED.speed_category,
  vedha = EXCLUDED.vedha,
  sign = EXCLUDED.sign,
  nakshatra_name = EXCLUDED.nakshatra_name,
  nakshatra_pada = EXCLUDED.nakshatra_pada,
  is_retro = EXCLUDED.is_retro,
  sign_lord = EXCLUDED.sign_lord,
  sign_lordship = EXCLUDED.sign_lordship,
  navamsa_sign = EXCLUDED.navamsa_sign,
  vargottama = EXCLUDED.vargottama,
  uday_bala = EXCLUDED.uday_bala,
  uchcha_bala = EXCLUDED.uchcha_bala,
  vakra_bala = EXCLUDED.vakra_bala,
  kshetra_bala = EXCLUDED.kshetra_bala,
  navamsha_bala = EXCLUDED.navamsha_bala,
  longitude_dms = EXCLUDED.longitude_dms,
  latitude_dms = EXCLUDED.latitude_dms,
  speed_long_dms = EXCLUDED.speed_long_dms;
