# syncdbdocs

syncdbdocs is a tool to help you keep your database documentation organized

## Problem to solve

While not in SQL standard, widely used databases have a way of commenting
tables and fields. This project intends to provide a way to generate
markdown and dbml documentation from those comments from the database
while at the same time making possible to update database comments from
the updates to the textual representation.

Said another way, you download current documentation to a markdown file
or dbml file that you can open with your text editor of choice and have a simple
view of the current schema and fields. You can edit markdown file and sync up
the new comments/editions back to the database.

Field order and table order will be preserved if you decide to reuse an
existing document. New fields or tables from the database will be appended.

## Formats

For now markdown and dbml are the two formats that this project will support.

Markdown will include all comments and some extra information (like data types),
while dbml is only provided to have a quick glance at the structure and
relationships of the data.

## Databases

postgres and mysql databases are supported. Right now only reading comments
and updating text/md files from database definitions is supported.

It should be easy to extend to other databases.

### PostgreSQL

- Read db definitions
- Update text/markdown from db
- Keep non-empty comments in the file if db has empty comments
- -Update db from text/markdown-

### MySQL

- Read db definitions
- Update text/markdown from db
- Keep non-empty comments in the file if db has empty comments
- -Update db from text/markdown-

Note for the future: when updating db from text files, we should be careful
since MySQL requires us to modify the whole column definition just to add
a comment. We can screw up data. Probably best way to proceed to avoid errors
is to create a temporary table, apply the definition there, and compare, if
things look the same, then we can safely apply the same alter table to the
original definition. Better to stay on the cautious side.
Use DESCRIBE, SHOW COLUMNS or SHOW CREATE TABLE.

## Project Status

Following there is a list of main features and whether or not they are supported.

Supported features:

- Support for postgres and mysql
- Generate markdown documentation
- Generate text documentation
- Generate DBMLish file (not standard, just to have a rough view of the structure)
- Update text & markdown from database without changing tables or field order


Missing features:

- Support for other databases: mssql, sqlite, ...
- Generate nicer HTML output (from text or database)
- Update database back
- Update markdown with new DB comments

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
