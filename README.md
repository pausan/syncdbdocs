# syncdbdocs

syncdbdocs is a tool to help you keep your database documented with a single
source of truth

## Problem to solve

While not in SQL standard, widely used databases have a way of commenting
tables and fields, althought it is usually very limited and cumbersome.

This project aims to provide a simple command to generate a textual representation
(in txt, markdown or dbml) of the database structure and be able to comment on it
and keep the documentation updated easily.

The original intend was to be able to sync back the documentation onto the
database, and while still possible and probably not hard to do, I've decided
that synching comments back to the database is not worth it due to the 
different constraints of some databases of limiting comment length and/or
comment availability altogether.

It is probably far more useful to preserve the documentation in textual form
and have a simple way of updating the document to include the new fields being
added.

Field order and table order will be preserved if you decide to reuse an
existing text or markdown file. New columns, tables and schemas from the database
will be appended in alphabetical order. Items that no longer exist will be marked
as deleted.

It is encouraged that you run this command automatically from your build after
migrating database schema and commit the resulting file.

## Installation

    $ go get github.com/pausan/syncdbdocs

or use docker

    $ docker pull pausan/syncdbdocs

## How to use

### Short version (TLDR)

First generate a document from scratch like this:

    $ syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -o pg_dbname.txt

Update documentation comments as you wish and then update database schema by
running the following command:

    $ syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -io pg_dbname.txt

You won't lose the editions of any of the existing schemas/tables/fields, but
if a schema/table/field is deleted from the database, it will be deleted from
the document (also renames).

It is encourage that you commit the documentation on your control version system
of choice.

### Docker Image

A tiny docker image is available so you can use it if you won't want to install
with go or you don't have a go compiler available.

You can use it like this:

    $ docker run pausan/syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -io pg_dbname.txt

### Long version

For now only commands to sync from the database to the file are provided
markdown and dbml documentation from those comments from the database
in order to keep a local textual representation of the schema.

Run the program once to generate the first version of the documentation, and,
in case there are no documented columns or tables or schema in the database
a document with the structure will be generated.

Please set DB_PASSWORD environment variable before running the command, and
select the database type with -t (or don't and it will try all drivers).

Example of initial import for all database types supported:

    $ syncdbdocs -t pg     -h 127.0.0.1 -u user -d dbname -o pg_dbname.txt
    $ syncdbdocs -t mysql  -h 127.0.0.1 -u user -d dbname -o mysql_dbname.txt
    $ syncdbdocs -t mssql  -h 127.0.0.1 -u user -d dbname -o mssql_dbname.txt
    $ syncdbdocs -t sqlite -h sqlite_file.db    -d dbname -o sqlite_dbname.txt

The command will write to stdout if no output file is provided.

To update the file after making some changes to the structure of the database
(eg adding tables or removing or renaming columns):

    $ syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -i pg_dbname.txt -o pg_dbname.txt

or:

    $ syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -io pg_dbname.txt

If you'd like to generate a markdown file or dbml instead use -format:

    $ syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -format markdown -o pg_dbname.md
    $ syncdbdocs -t pg -h 127.0.0.1 -u user -d dbname -format dbml -o pg_dbname.dbml

If you want to check out more parameters, just run with -h or -help.

## Formats

Plain **text** files, **markdown** and **dbml** are the supported formats.

Markdown and text files include all comments and some extra information (like data types),
while dbml is only provided to have a quick glance at the structure of the data.

## Databases

postgres, mysql and mssql are supported. Right now only reading comments
and updating text/md files from database definitions is supported.

It should be easy to extend to other databases.

### PostgreSQL

- Read db definitions
- Update text/markdown from db
- Keep non-empty comments in the file if db has empty comments
- Tested with postgres 9.x, 10.x, 11.x and 12.x

### MySQL

- Read db definitions
- Update text/markdown from db
- Keep non-empty comments in the file if db has empty comments
- Tested with mysql v8.x

Note for the future: when updating db from text files, we should be careful
since MySQL requires us to modify the whole column definition just to add
a comment. We can screw up data. Probably best way to proceed to avoid errors
is to create a temporary table, apply the definition there, and compare, if
things look the same, then we can safely apply the same alter table to the
original definition. Better to stay on the cautious side.
Use DESCRIBE, SHOW COLUMNS or SHOW CREATE TABLE.

### MS SQL Server

- Read db definitions
- Update text/markdown from db
- Keep non-empty comments in the file if db has empty comments
- Tested with sql server 2017 and 2019

### SQLite

This database does not support comments, but this tool supports pulling the
database structure from it to update a comment file.

- Read db definitions
- Update text/markdown from db
- Tested with sqlite 3.x

## Project Status

Following there is a list of main features and whether or not they are supported.

Supported features:

- Support for postgres, mysql, mssql and sqlite
- Generate/update markdown documentation
- Generate/update text documentation
- Generate DBMLish file (not standard, just to have a rough view of the structure)
- Update text & markdown from database without changing tables or field order

Missing features:

- Update database back from file comments
- Support for other databases: oracle, ...
- Generate nicer HTML output (from text or database)
- Detect primary keys, indexes, triggers or functions

## License

This software uses GPL v2 license, you can read it fully in LICENSE file.

    Copyright (C) 2021 Pau Sanchez

    This program is free software; you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation; either version 2 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License along
    with this program; if not, write to the Free Software Foundation, Inc.,
    51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
