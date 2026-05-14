-- Adminer 5.4.2 PostgreSQL 18.3 dump


CREATE TYPE "exchange_type" AS ENUM ('BSE', 'NSE');


CREATE TYPE "planet_type" AS ENUM ('Sun', 'Moon', 'Mercury', 'Venus', 'Mars', 'Jupiter', 'Saturn', 'Uranus', 'Neptune', 'Pluto', 'Rahu', 'Ketu');


CREATE TYPE "vedha_type" AS ENUM ('left', 'right', 'front', 'no', 'n/a');


CREATE TYPE "rel_type" AS ENUM ('Friend', 'Neutral', 'Enemy', 'Self');


CREATE TYPE "week_day_type" AS ENUM ('Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday');


CREATE TYPE "sign_type" AS ENUM ('Aries', 'Taurus', 'Gemini', 'Cancer', 'Leo', 'Virgo', 'Libra', 'Scorpio', 'Sagittarius', 'Capricorn', 'Aquarius', 'Pisces');


CREATE TYPE "nakshatra_type" AS ENUM ('Ashwini', 'Bharani', 'Krittika', 'Rohini', 'Mrigashirsha', 'Ardra', 'Punarvasu', 'Pushya', 'Ashlesha', 'Magha', 'Purva Phalguni', 'Uttara Phalguni', 'Hasta', 'Chitra', 'Swati', 'Vishakha', 'Anuradha', 'Jyestha', 'Moola', 'Purva Ashadha', 'Uttara Ashadha', 'Abhijit', 'Shravana', 'Dhanishtha', 'Shatabhisha', 'Purva Bhadrapada', 'Uttara Bhadrapada', 'Revati');

CREATE TYPE "speed_type" AS ENUM ('kutil', 'ati-vakra', 'vakra', 'ati-mand', 'mand', 'madhyam', 'sama', 'sheeghra', 'ati-sheeghra', 'n/a');

CREATE TABLE "tbl_ohlc" (
    "day" timestamptz NOT NULL,
    "ticker_id" smallint NOT NULL,
    "o" double precision,
    "h" double precision,
    "l" double precision,
    "c" double precision,
    "v" double precision,
    CONSTRAINT "tbl_ohlc_pkey" PRIMARY KEY ("day", "ticker_id")
);



CREATE TABLE "tbl_panchang" (
    "time" timestamptz NOT NULL,
    "tithy" smallint NOT NULL,
    "nakshatra" nakshatra_type NOT NULL,
    "weekday" week_day_type NOT NULL,
    CONSTRAINT "tbl_panchang_pkey" PRIMARY KEY ("time")
);

TRUNCATE "tbl_panchang";

CREATE TABLE "tbl_planet_positions" (
    "observation_time" timestamptz NOT NULL,
    "planet_name" planet_type NOT NULL,
    "longitude" double precision NOT NULL,
    "latitude" double precision NOT NULL,
    "distance" double precision NOT NULL,
    "speed_long" double precision NOT NULL,
    "speed_lat" double precision NOT NULL,
    "speed_dist" double precision NOT NULL,
    "speed_category" speed_type NOT NULL,
    "vedha" vedha_type NOT NULL,
    "sign" sign_type,
    "nakshatra_name" nakshatra_type,
    "nakshatra_pada" smallint,
    "is_retro" boolean DEFAULT false,
    "sign_lord" planet_type,
    "sign_lordship" rel_type,
    "navamsa_sign" sign_type,
    "vargottama" boolean DEFAULT false,
    "uday_bala" double precision,
    "uchcha_bala" double precision,
    "vakra_bala" double precision,
    "kshetra_bala" double precision,
    "navamsha_bala" double precision,
    "longitude_dms" jsonb,
    "latitude_dms" jsonb,
    "speed_long_dms" jsonb,
    CONSTRAINT "tbl_planet_positions_pkey" PRIMARY KEY ("observation_time", "planet_name")
);

TRUNCATE "tbl_planet_positions";

CREATE TABLE "tbl_ticker" (
    "id" smallint NOT NULL,
    "name" character varying(30) NOT NULL,
    "desc" text,
    "full_name" text,
    "exchange" exchange_type NOT NULL,
    CONSTRAINT "tbl_ticker_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "check_id_limit" CHECK ((id <= 10000))
);

TRUNCATE "tbl_ticker";
INSERT INTO "tbl_ticker" ("id", "name", "desc", "full_name", "exchange") VALUES
(2,	'TCS',	'IT Services and Consulting',	'Tata Consultancy Services Limited',	'BSE'),
(1,	'NIFTY50',	'nifty 50',	'Nifty 50',	'NSE');

-- 2026-05-13 21:54:33 UTC
