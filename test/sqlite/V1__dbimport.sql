-- -- sqlite does not seem to have any kind of comments

-- ------------------------------------------------------------------------------
-- multiple_types
-- ------------------------------------------------------------------------------
CREATE TABLE multiple_types (
  id INTEGER PRIMARY KEY AUTOINCREMENT,

  _int INT,
  _integer INTEGER DEFAULT 32,
  _tinyint TINYINT,
  _smallint SMALLINT,
  _mediumint MEDIUMINT,
  _bigint BIGINT,
  _ubigint UNSIGNED BIG INT,
  _int2 INT2,
  _int8 INT8,

  _character CHARACTER(20),
  _varchar VARCHAR(255),
  _varchar2 VARYING CHARACTER(25),
  _nchar NCHAR(55),
  _natchar NATIVE CHARACTER(70),
  _nvarchar NVARCHAR(100),
  _text TEXT,
  _clob CLOB,

  _blob BLOB,
  _real REAL,
  _double DOUBLE,
  _double_precision DOUBLE PRECISION,
  _float FLOAT,

  _numeric NUMERIC,
  _decimal DECIMAL(10,5),
  _boolean BOOLEAN,
  _date DATE,
  _datetime DATETIME
);

-- ------------------------------------------------------------------------------
-- user
-- ------------------------------------------------------------------------------
CREATE TABLE user (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  full_name    VARCHAR(128) DEFAULT NULL,
  email        VARCHAR(128) UNIQUE NOT NULL,
  password     VARCHAR(256) NOT NULL,
  access       TEXT NOT NULL DEFAULT 'NONE',
  language     CHAR(2) DEFAULT NULL,
  country_code CHAR(2) NOT NULL,
  created_date TIMESTAMP NOT NULL,
  updated_date TIMESTAMP NOT NULL
);

