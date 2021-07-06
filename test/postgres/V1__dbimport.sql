
CREATE EXTENSION IF NOT EXISTS  pgcrypto;

CREATE SCHEMA syncdbtest;

COMMENT ON DATABASE dbtest IS 'Hey!! This is a comment about the database we are documenting, it should appear the first one, and should logically wrap to whatever max line width you specify in syncdbdocs command line.';
COMMENT ON SCHEMA syncdbtest IS 'Let''s see how this comment about the schema works out';

--------------------------------------------------------------------------------
-- uint2: custom unsigned int2 value
--------------------------------------------------------------------------------
CREATE DOMAIN uint2 AS int4
   CHECK(VALUE >= 0 AND VALUE < 65536);

--------------------------------------------------------------------------------
-- syncdbtest.access_level
--------------------------------------------------------------------------------
CREATE TYPE syncdbtest.access_level AS ENUM('NONE', 'VIEW', 'EDIT', 'ADMIN');

--------------------------------------------------------------------------------
-- syncdbtest.syncUpdateDate trigger
--------------------------------------------------------------------------------
CREATE FUNCTION syncdbtest.syncUpdatedDate() RETURNS trigger AS $$
BEGIN
  NEW.updated_date := (NOW() AT TIME ZONE 'UTC');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

--------------------------------------------------------------------------------
-- syncdbtest.multiple_types
--------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS syncdbtest.multiple_types (
  _access_level   syncdbtest.access_level NOT NULL,
  _varchar16      CHARACTER VARYING(64) NOT NULL DEFAULT '_varchar16 value',
  _varchar64      VARCHAR(64) NOT NULL DEFAULT '_varchar64 value',
  _char2          CHAR(2),
  _char16         CHAR(16),

  _smallintcheck  SMALLINT CHECK (_smallintcheck > 1234),
  _uint2          uint2,

  _bigint         BIGINT,
  _bigserial      BIGSERIAL,
  _bit            BIT,
  _boolean        BOOLEAN,
  _box            BOX,
  _bytea          BYTEA,
  _character      CHARACTER,
  _cidr           CIDR,
  _circle         CIRCLE,
  _date           DATE,
  _double         FLOAT8,
  _inet           INET,
  _integer        INTEGER,
  _interval       INTERVAL,
  _json           JSON,
  _jsonb          JSONB,
  _line           LINE,
  _lseg           LSEG,
  _macaddr        MACADDR,
  _money          MONEY,
  _numeric        NUMERIC,
  _path           PATH,
  _pg_lsn         PG_LSN,
  _point          POINT,
  _polygon        POLYGON,
  _real           REAL,
  _smallint       SMALLINT,
  _smallserial    SMALLSERIAL,
  _serial         SERIAL,
  _text           TEXT DEFAULT NULL,
  _time           TIME,
  _timestamp      TIMESTAMP,
  _tsquery        TSQUERY,
  _tsvector       TSVECTOR,
  _txid_snapshot  TXID_SNAPSHOT,
  _uuid           UUID PRIMARY KEY NOT NULL,
  _xml            XML
);

--------------------------------------------------------------------------------
-- syncdbtest.user
--------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS syncdbtest.user (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  full_name    VARCHAR(128) DEFAULT NULL,
  email        VARCHAR(128) UNIQUE NOT NULL,

  password     VARCHAR(256) NOT NULL,
  access       syncdbtest.access_level NOT NULL DEFAULT 'NONE',

  language     CHAR(2) DEFAULT NULL,
  country_code CHAR(2) NOT NULL,

  created_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
  updated_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TRIGGER syncUpdatedDate
  BEFORE UPDATE ON syncdbtest.user
  FOR EACH ROW EXECUTE PROCEDURE syncdbtest.syncUpdatedDate();

COMMENT ON TABLE syncdbtest.user IS
  'This is the test comment that we are going to use for the user table, we can make it simpler, but this is long because we also want to test how good the algorithm of word-wrap works sorting things out; I believe it will work well, but we will see.';

COMMENT ON COLUMN syncdbtest.user.language IS
  'Language represents a ISO-639-2 standard value';

COMMENT ON COLUMN syncdbtest.user.country_code IS
  'Country code represents a ISO-3166 alpha-2 value. Should not be NULL.';

COMMENT ON COLUMN syncdbtest.user.access IS
  'Access level that this user has in the current system';

COMMENT ON COLUMN syncdbtest.user.email IS
  'As you have figured out, this is the email address of the user';

COMMENT ON COLUMN syncdbtest.user.password IS
  'Password *** _ ## \\ \\`{}[]<>()#*+-_.!| **markdown** escape check';
