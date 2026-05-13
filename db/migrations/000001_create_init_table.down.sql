
BEGIN;

TRUNCATE TABLE tbl_ticker;
DROP TABLE IF EXISTS tbl_ticker;
DROP TYPE exchange_type;
DROP EXTENSION IF EXISTS timescaledb;

COMMIT;
