# dbtest (PostgreSQL)

Hey\!\! This is a comment about the database we are documenting, it should
appear the first one, and should logically wrap to whatever max line width you
specify in syncdbdocs command line.

## public

standard public schema

### flyway\_schema\_history

- checksum [int4?]

- description [varchar200]

- execution\_time [int4]

- installed\_by [varchar100]

- installed\_on [timestamp]

- installed\_rank [int4]

- script [varchar1000]

- success [bool]

- type [varchar20]

- version [varchar50?]

## syncdbtest

Let's see how this comment about the schema works out

### multiple\_types

- \_access\_level [access\_level]

- \_bigint [int8?]

- \_bigserial [int8]

- \_bit [bit1?]

- \_boolean [bool?]

- \_box [box?]

- \_bytea [bytea?]

- \_char16 [bpchar16?]

- \_char2 [bpchar2?]

- \_character [bpchar1?]

- \_cidr [cidr?]

- \_circle [circle?]

- \_date [date?]

- \_double [float8?]

- \_inet [inet?]

- \_integer [int4?]

- \_interval [interval?]

- \_json [json?]

- \_jsonb [jsonb?]

- \_line [line?]

- \_lseg [lseg?]

- \_macaddr [macaddr?]

- \_money [money?]

- \_numeric [numeric?]

- \_path [path?]

- \_pg\_lsn [pg\_lsn?]

- \_point [point?]

- \_polygon [polygon?]

- \_real [float4?]

- \_serial [int4]

- \_smallint [int2?]

- \_smallintcheck [int2?]

- \_smallserial [int2]

- \_text [text?]

- \_time [time?]

- \_timestamp [timestamp?]

- \_tsquery [tsquery?]

- \_tsvector [tsvector?]

- \_txid\_snapshot [txid\_snapshot?]

- \_uint2 [int4?]

- \_uuid [uuid]

- \_varchar16 [varchar64]

- \_varchar64 [varchar64]

- \_xml [xml?]

### user

This is the test comment that we are going to use for the user table, we can
make it simpler, but this is long because we also want to test how good the
algorithm of word\-wrap works sorting things out; I believe it will work well,
but we will see.

- access [access\_level]

  Access level that this user has in the current system

- country\_code [bpchar2]

  Country code represents a ISO\-3166 alpha\-2 value. Should not be NULL.

- created\_date [timestamp]

- email [varchar128]

  As you have figured out, this is the email address of the user

- full\_name [varchar128?]

- id [uuid]

- language [bpchar2?]

  Language represents a ISO\-639\-2 standard value

- password [varchar256]

  Password \*\*\* \_ \#\# \\\\ \\\\\`\{\}\[\]\<\>\(\)\#\*\+\-\_.\!\|
  \*\*markdown\*\* escape check

- updated\_date [timestamp]

