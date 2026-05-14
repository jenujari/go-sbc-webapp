BEGIN;

ALTER TABLE tbl_ohlc
DROP CONSTRAINT tbl_ohlc_pkey;

DROP TABLE IF EXISTS tbl_ohlc;

DROP TABLE IF EXISTS tbl_planet_positions;
DROP TABLE IF EXISTS tbl_panchang;

DROP TYPE IF EXISTS planet_type;
DROP TYPE IF EXISTS vedha_type;
DROP TYPE IF EXISTS sign_type;
DROP TYPE IF EXISTS nakshatra_type;
DROP TYPE IF EXISTS rel_type;
DROP TYPE IF EXISTS week_day_type;
DROP TYPE IF EXISTS speed_type;


COMMIT;
