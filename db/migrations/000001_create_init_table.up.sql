BEGIN;

CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TYPE exchange_type AS ENUM ('BSE', 'NSE');

CREATE TABLE tbl_ticker (
    id SMALLINT PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    "desc" TEXT,
    full_name TEXT,
    exchange exchange_type NOT NULL,

    -- Constraint to ensure id does not exceed 10,000
    CONSTRAINT check_id_limit CHECK (id <= 10000)
);

INSERT INTO tbl_ticker (id, name, "desc", full_name, exchange)
VALUES
    (1, 'RELIANCE', 'Energy and Telecom giant', 'Reliance Industries Limited', 'NSE'),
    (2, 'TCS', 'IT Services and Consulting', 'Tata Consultancy Services Limited', 'BSE');

COMMIT;
