
-- mysql does not seem to have any database level comment

-- ------------------------------------------------------------------------------
-- multiple_types
-- ------------------------------------------------------------------------------
CREATE TABLE multiple_types (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
 _bit BIT(8),
 _bool bool,
 _smallint SMALLINT,
 _mediumint MEDIUMINT,
 _int INT,
 _bigint BIGINT,
 _float FLOAT,
 _double DOUBLE,
 _decimal DECIMAL(4,2),
 _enum ENUM("a", "b", "c"),
 _char2 CHAR(2),
 _varchar16 VARCHAR(16),
 _varchar64 VARCHAR(64),
 _binary255 BINARY(255),
 _varbinary255 VARBINARY(255),
 _text TEXT,
 _blob BLOB,
 _blob_1k BLOB(1024),
 _tinyblob TINYBLOB,
 _tinytext TINYTEXT,
 _set SET("a", "b", "c", "d")
);

-- ------------------------------------------------------------------------------
-- user
-- ------------------------------------------------------------------------------
CREATE TABLE user (
  id           BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID(), TRUE)),

  full_name    VARCHAR(128) DEFAULT NULL,

  email        VARCHAR(128) UNIQUE NOT NULL
               COMMENT 'As you have figured out, this is the email address of the user',

  password     VARCHAR(256) NOT NULL
               COMMENT 'Password *** _ ## \\ \\`{}[]<>()#*+-_.!| **markdown** escape check',

  access       ENUM('NONE', 'READ', 'EDIT', 'ADMIN') NOT NULL DEFAULT 'NONE'
               COMMENT 'Access level that this user has in the current system',

  language     CHAR(2) DEFAULT NULL
               COMMENT 'Language represents a ISO-639-2 standard value',

  country_code CHAR(2) NOT NULL
               COMMENT 'Country code represents a ISO-3166 alpha-2 value. Should not be NULL.',

  created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  updated_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
      ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE user COMMENT
  'This is the test comment that we are going to use for the user table, we can make it simpler, but this is long because we also want to test how good the algorithm of word-wrap works sorting things out; I believe it will work well, but we will see.';

