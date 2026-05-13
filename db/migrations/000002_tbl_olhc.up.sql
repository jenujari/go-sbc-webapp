BEGIN;

CREATE TABLE tbl_ohlc (
    day TIMESTAMPTZ NOT NULL,
    ticker_id SMALLINT NOT NULL,
    o DOUBLE PRECISION,
    h DOUBLE PRECISION,
    l DOUBLE PRECISION,
    c DOUBLE PRECISION,
    v DOUBLE PRECISION,
    CONSTRAINT fk_ticker
        FOREIGN KEY (ticker_id)
        REFERENCES tbl_ticker (id)
);

ALTER TABLE tbl_ohlc
ADD PRIMARY KEY (day, ticker_id);

CREATE TYPE planet_type AS ENUM ('Sun', 'Moon', 'Mercury', 'Venus', 'Mars', 'Jupiter', 'Saturn', 'Uranus', 'Neptune', 'Pluto', 'Rahu', 'Ketu');
CREATE TYPE vedha_type AS ENUM ('left', 'right', 'front', 'no', 'n/a');
CREATE TYPE rel_type AS ENUM ('Friend', 'Neutral', 'Enemy', 'Self');
CREATE TYPE week_day_type AS ENUM ('Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday');
CREATE TYPE sign_type AS ENUM ('Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo', 'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces');
CREATE TYPE nakshatra_type AS ENUM ('Ashwini', 'Bharani', 'Krittika', 'Rohini', 'Mrigashirsha', 'Ardra', 'Punarvasu', 'Pushya', 'Ashlesha', 'Magha', 'Purva Phalguni', 'Uttara Phalguni', 'Hasta', 'Chitra', 'Swati', 'Vishakha', 'Anuradha', 'Jyestha', 'Moola', 'Purva Ashadha', 'Uttara Ashadha', 'Abhijit', 'Shravana', 'Dhanishtha', 'Shatabhisha', 'Purva Bhadrapada', 'Uttara Bhadrapada', 'Revati');
CREATE TYPE speed_type AS ENUM ('kutil', 'ati-vakra', 'vakra', 'ati-mand', 'mand', 'madhyam', 'sama', 'sheeghra', 'ati-sheeghra','n/a');


CREATE TABLE IF NOT EXISTS tbl_planet_positions (
    -- Time-series anchor (Required for Hypertable conversion later)
    observation_time TIMESTAMPTZ NOT NULL,

    -- Planet Identification
    planet_name      planet_type NOT NULL,

    -- Coordinate System (Double Precision for high accuracy)
    longitude        DOUBLE PRECISION NOT NULL,
    latitude         DOUBLE PRECISION NOT NULL,
    distance         DOUBLE PRECISION NOT NULL,

    -- Velocity Data
    speed_long       DOUBLE PRECISION NOT NULL,
    speed_lat        DOUBLE PRECISION NOT NULL,
    speed_dist       DOUBLE PRECISION NOT NULL,
    speed_category   speed_type NOT NULL,

    -- Vedic / Astrological Context
    vedha            vedha_type NOT NULL,
    sign             sign_type,
    nakshatra_name   nakshatra_type,
    nakshatra_pada   SMALLINT,
    is_retro         BOOLEAN DEFAULT FALSE,
    sign_lord        planet_type,
    sign_lordship    rel_type,
    navamsa_sign     sign_type,
    vargottama       BOOLEAN DEFAULT FALSE,

    -- Calculated Bala (Strengths)
    uday_bala        DOUBLE PRECISION,
    uchcha_bala      DOUBLE PRECISION,
    vakra_bala       DOUBLE PRECISION,
    kshetra_bala     DOUBLE PRECISION,
    navamsha_bala    DOUBLE PRECISION,

    -- DMS Data (Stored as JSONB for flexibility, or could be expanded to columns)
    longitude_dms    JSONB,
    latitude_dms     JSONB,
    speed_long_dms   JSONB,

    -- Composite Primary Key for TimescaleDB compatibility
    PRIMARY KEY (observation_time, planet_name)
);

CREATE TABLE IF NOT EXISTS tbl_panchang (
    time        TIMESTAMPTZ NOT NULL,
    tithy       SMALLINT    NOT NULL,
    nakshatra   nakshatra_type NOT NULL,
    weekday     week_day_type NOT NULL,

    -- Primary Key must include the partitioning column (time)
    PRIMARY KEY (time)
);

COMMIT;
