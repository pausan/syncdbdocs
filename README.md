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

Only postgres database is supported. Hopefully it is built in a way that can be
easily extended to other databases as well.

## Project Status

Following there is a list of main features and whether or not they are supported.

Supported features:

- Support for postgres
- Generate markdown documentation

Missing features:

- Support for mysql, mssql, sqlite, ...
- Generate DBML file
- Read markdown updates
- Update database back
- Update markdown with new DB changes

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
