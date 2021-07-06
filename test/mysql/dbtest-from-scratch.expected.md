# dbtest (MySQL)

### flyway\_schema\_history

- checksum [int?]

- description [varchar\(200\)]

- execution\_time [int]

- installed\_by [varchar\(100\)]

- installed\_on [timestamp]

- installed\_rank [int]

- script [varchar\(1000\)]

- success [tinyint\(1\)]

- type [varchar\(20\)]

- version [varchar\(50\)?]

### multiple\_types

- \_bigint [bigint?]

- \_binary255 [binary\(255\)?]

- \_bit [bit\(8\)?]

- \_blob [blob?]

- \_blob\_1k [blob?]

- \_bool [tinyint\(1\)?]

- \_char2 [char\(2\)?]

- \_decimal [decimal\(4,2\)?]

- \_double [double?]

- \_enum [enum\('a','b','c'\)?]

- \_float [float?]

- \_int [int?]

- \_mediumint [mediumint?]

- \_set [set\('a','b','c','d'\)?]

- \_smallint [smallint?]

- \_text [text?]

- \_tinyblob [tinyblob?]

- \_tinytext [tinytext?]

- \_varbinary255 [varbinary\(255\)?]

- \_varchar16 [varchar\(16\)?]

- \_varchar64 [varchar\(64\)?]

- id [int unsigned]

### user

This is the test comment that we are going to use for the user table, we can
make it simpler, but this is long because we also want to test how good the
algorithm of word\-wrap works sorting things out; I believe it will work well,
but we will see.

- access [enum\('NONE','READ','EDIT','ADMIN'\)]

  Access level that this user has in the current system

- country\_code [char\(2\)]

  Country code represents a ISO\-3166 alpha\-2 value. Should not be NULL.

- created\_date [timestamp]

- email [varchar\(128\)]

  As you have figured out, this is the email address of the user

- full\_name [varchar\(128\)?]

- id [binary\(16\)]

- language [char\(2\)?]

  Language represents a ISO\-639\-2 standard value

- password [varchar\(256\)]

  Password \*\*\* \_ \#\# \\ \\\`\{\}\[\]\<\>\(\)\#\*\+\-\_.\!\|
  \*\*markdown\*\* escape check

- updated\_date [timestamp]

