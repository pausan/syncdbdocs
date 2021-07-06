# dbtest (PostgreSQL)

Hey\!\! This is a comment about the database we are documenting, it should
appear the first one, and should logically wrap to whatever max line width you
specify in syncdbdocs command line.

## syncdbtest

Let's see how this comment about the schema works out

### multiple\_types

- \_cidr [cidr?]

- \_tsquery [tsquery?]

- \_interval [interval?]

- \_char16 [bpchar16?]

- \_bytea [bytea?]

- \_inet [inet?]

- \_text [text?]

- \_time [time?]

- \_polygon [polygon?]

- \_line [line?]

- \_integer [int4?]

- \_money [money?]

- \_lseg [lseg?]

- \_macaddr [macaddr?]

- \_date [date?]

- \_real [float4?]

- \_point [point?]

- \_boolean [bool?]

- \_timestamp [timestamp?]

- \_varchar16 [varchar64]

- \_smallint [int2?]

- \_bit [bit1?]

- \_box [box?]

- \_double [float8?]

- \_serial [int4]

- \_char2 [bpchar2?]

- \_jsonb [jsonb?]

- \_xml [xml?]

- \_circle [circle?]

- \_access\_level [access\_level]

- \_json [json?]

- \_character [bpchar1?]

- \_smallserial [int2]

- \_uuid [uuid]

- \_varchar64 [varchar64]

- \_uint2 [int4?]

- \_bigserial [int8]

- \_numeric [numeric?]

- \_pg\_lsn [pg\_lsn?]

- \_bigint [int8?]

- \_smallintcheck [int2?]

- \_path [path?]

- \_tsvector [tsvector?]

- \_txid\_snapshot [txid\_snapshot?]

### user

This is the test comment that we are going to use for the user table, we can
make it simpler, but this is long because we also want to test how good the
algorithm of word\-wrap works sorting things out; I believe it will work well,
but we will see.

- full\_name [varchar128?]

- access [access\_level]

  Access level that this user has in the current system

- language [bpchar2?]

  Language represents a ISO\-639\-2 standard value

- created\_date [timestamp]

- email [varchar128]

  As you have figured out, this is the email address of the user

- country\_code [bpchar2]

  Country code represents a ISO\-3166 alpha\-2 value. Should not be NULL.

- updated\_date [timestamp]

- password [varchar256]

  Password \*\*\* \_ \#\# \\\\ \\\\\`\{\}\[\]\<\>\(\)\#\*\+\-\_.\!\|
  \*\*markdown\*\* escape check

- id [uuid]

## public

standard public schema

### flyway\_schema\_history

- checksum [int4?]

- installed\_on [timestamp]

- description [varchar200]

- execution\_time [int4]

- installed\_by [varchar100]

- type [varchar20]

- installed\_rank [int4]

- success [bool]

- version [varchar50?]

- script [varchar1000]

