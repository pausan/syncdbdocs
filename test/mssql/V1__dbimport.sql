
USE [master]
GO

CREATE DATABASE [dbtest]
GO

USE [dbtest]
GO

CREATE SCHEMA [syncdbtest]
GO

CREATE TABLE syncdbtest.[user](
  id           INT,

  full_name    VARCHAR(128) DEFAULT NULL,
  email        VARCHAR(128) UNIQUE NOT NULL,

  password     VARCHAR(256) NOT NULL,
  access       VARCHAR(10) NOT NULL DEFAULT 'NONE' CHECK (access IN('NONE', 'READ', 'EDIT', 'ADMIN')),

  language     CHAR(2) DEFAULT NULL,
  country_code CHAR(2) NOT NULL,

  created_date DATETIME NOT NULL,
  updated_date DATETIME NOT NULL
)
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'Hey!! This is a comment about the database we are documenting, it should appear the first one, and should logically wrap to whatever max line width you specify in syncdbdocs command line.'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'Let''s see how this comment about the schema works out.',
  'schema', N'syncdbtest'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'This is the test comment that we are going to use for the user table, we can make it simpler, but this is long because we also want to test how good the algorithm of word-wrap works sorting things out; I believe it will work well, but we will see.',
  'schema', N'syncdbtest',
  'table', N'user'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'Language represents a ISO-639-2 standard value',
  'schema', N'syncdbtest',
  'table', N'user',
  'column', N'language'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'Country code represents a ISO-3166 alpha-2 value. Should not be NULL.',
  'schema', N'syncdbtest',
  'table', N'user',
  'column', N'country_code'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'Access level that this user has in the current system',
  'schema', N'syncdbtest',
  'table', N'user',
  'column', N'access'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'As you have figured out, this is the email address of the user',
  'schema', N'syncdbtest',
  'table', N'user',
  'column', N'email'
GO

EXEC dbtest.sys.sp_addextendedproperty
  'MS_Description', N'Password *** _ ## \\ \\`{}[]<>()#*+-_.!| **markdown** escape check',
  'schema', N'syncdbtest',
  'table', N'user',
  'column', N'password'
