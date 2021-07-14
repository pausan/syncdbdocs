#
# Creates some databases, fill them with data, and then run syncdbdocs apps on
# them to check if works as expected
#
RANDVAR      ?= default
NETWORK_NAME ?= sdd_network_$(RANDVAR)

SYNCDBDOCS_IMAGE ?= syncdbdocs

PG_CONTAINER        ?= sdd_pg_$(RANDVAR)
MYSQL_CONTAINER     ?= sdd_mysql_$(RANDVAR)
MSSQL_CONTAINER     ?= sdd_mssql_$(RANDVAR)

PG_IMAGE_VERSION    ?= latest
MYSQL_IMAGE_VERSION ?= latest
MSSQL_IMAGE_VERSION ?= 2019-latest

DB_NAME = dbtest
DB_USER = root
DB_PASS = pass

PG_PORT    = 5432
MYSQL_PORT = 3306

MSSQL_PORT = 1433
MSSQL_USER = sa
MSSQL_PASS = _asdfASDF123

SQLITE_FILE = testdb.db


.PHONY: all test clean

all:   | up build test clean
clean: | down

up:    | net-up db-up migrate
down:  | db-down net-down

net-up:
	docker network create $(NETWORK_NAME) || true

net-down:
	docker network rm $(NETWORK_NAME) || true

db-up: | pg-up mysql-up mssql-up

pg-up:
	docker run -d --rm \
	  --name $(PG_CONTAINER) \
	  --network $(NETWORK_NAME) \
	  -e POSTGRES_DB=$(DB_NAME) \
	  -e POSTGRES_USER=$(DB_USER) \
	  -e POSTGRES_PASSWORD=$(DB_PASS) \
	  postgres:$(PG_IMAGE_VERSION)

mysql-up:
	docker run -d --rm \
	  --name $(MYSQL_CONTAINER) \
	  --network $(NETWORK_NAME) \
	  -e MYSQL_DATABASE=$(DB_NAME) \
	  -e MYSQL_ROOT_PASSWORD=$(DB_PASS) \
	  mysql:$(MYSQL_IMAGE_VERSION)

mssql-up:
	docker run -d --rm \
	  --name $(MSSQL_CONTAINER) \
	  --network $(NETWORK_NAME) \
	  -e ACCEPT_EULA=Y \
	  -e SA_PASSWORD=$(MSSQL_PASS) \
	  mcr.microsoft.com/mssql/server:$(MSSQL_IMAGE_VERSION)

migrate: | migrate-pg migrate-mysql migrate-mssql

migrate-pg:
	docker run --rm \
	  --network $(NETWORK_NAME) \
		-v $(PWD)/test/postgres:/flyway/sql \
		flyway/flyway:7.11 \
		-url=jdbc:postgresql://$(PG_CONTAINER):$(PG_PORT)/$(DB_NAME) \
		-user=$(DB_USER) \
		-password=$(DB_PASS) \
		-connectRetries=5 \
		migrate

migrate-mysql:
	docker run --rm \
	  --network $(NETWORK_NAME) \
		-v $(PWD)/test/mysql:/flyway/sql \
		flyway/flyway:7.11 \
		-url=jdbc:mysql://$(MYSQL_CONTAINER):$(MYSQL_PORT)/$(DB_NAME) \
		-user=$(DB_USER) \
		-password=$(DB_PASS) \
		-connectRetries=5 \
		migrate

migrate-mssql:
	docker run --rm \
	  --network $(NETWORK_NAME) \
		-v $(PWD)/test/mssql:/flyway/sql \
		flyway/flyway:7.11 \
		-mixed=true \
		"-url=jdbc:sqlserver://$(MSSQL_CONTAINER):$(MSSQL_PORT);databaseName=master" \
		-user=$(MSSQL_USER) \
		-password=$(MSSQL_PASS) \
		-connectRetries=5 \
		migrate

migrate-sqlite:
	rm -f $(PWD)/test/sqlite/$(SQLITE_FILE)
	touch $(PWD)/test/sqlite/$(SQLITE_FILE)
	docker run --rm \
		-v $(PWD)/test/sqlite:/db:rw \
		debian:stable \
		bash -c " \
		apt update \
		&& apt install sqlite3 -y  \
		&& sqlite3 /db/$(SQLITE_FILE) < /db/V1__dbimport.sql \
		"

db-down:
	docker stop $(PG_CONTAINER) || true
	docker stop $(MYSQL_CONTAINER) || true
	docker stop $(MSSQL_CONTAINER) || true

build:
	docker build . --tag $(SYNCDBDOCS_IMAGE)


PG_RUN_SYNCDBDOCS = docker run --rm \
	--network $(NETWORK_NAME) \
	-e DB_PASSWORD=$(DB_PASS) \
	-v $(PWD)/test/postgres:/tmp/testpg/:ro \
	$(SYNCDBDOCS_IMAGE) \
	-h $(PG_CONTAINER) \
	-p $(PG_PORT) \
	-u $(DB_USER) \
	-d $(DB_NAME)

test-pg:
	$(PG_RUN_SYNCDBDOCS) -format=md > /tmp/dbtest.result.md
	$(PG_RUN_SYNCDBDOCS) -format=text > /tmp/dbtest.result.txt
	diff $(PWD)/test/postgres/dbtest-from-scratch.expected.md /tmp/dbtest.result.md || (echo "PG Test001.md failed" && false)
	diff $(PWD)/test/postgres/dbtest-from-scratch.expected.txt /tmp/dbtest.result.txt || (echo "PG Test001.txt failed" && false)

	$(PG_RUN_SYNCDBDOCS) -db-comments-first -i /tmp/testpg/dbtest-preserve-order.input > /tmp/dbtest.result
	diff $(PWD)/test/postgres/dbtest-preserve-order.expected.txt /tmp/dbtest.result || (echo "PG Test002 failed" && false)

	$(PG_RUN_SYNCDBDOCS) -i /tmp/testpg/dbtest-preserve-order.input > /tmp/dbtest.result
	diff $(PWD)/test/postgres/dbtest-preserve-file-comments.expected.txt /tmp/dbtest.result || (echo "PG Test003 failed" && false)

MYSQL_RUN_SYNCDBDOCS = docker run --rm \
	--network $(NETWORK_NAME) \
	-e DB_PASSWORD=$(DB_PASS) \
	-v $(PWD)/test/mysql:/tmp/testmysql/:ro \
	$(SYNCDBDOCS_IMAGE) \
	-h $(MYSQL_CONTAINER) \
	-p $(MYSQL_PORT) \
	-u $(DB_USER) \
	-d $(DB_NAME)

test-mysql:
	$(MYSQL_RUN_SYNCDBDOCS) -format=md > /tmp/dbtest.result.md
	$(MYSQL_RUN_SYNCDBDOCS) -format=text > /tmp/dbtest.result.txt
	diff $(PWD)/test/mysql/dbtest-from-scratch.expected.md /tmp/dbtest.result.md || (echo "MYSQL Test001.md failed" && false)
	diff $(PWD)/test/mysql/dbtest-from-scratch.expected.txt /tmp/dbtest.result.txt || (echo "MYSQL Test001.txt failed" && false)

MSSQL_RUN_SYNCDBDOCS = docker run --rm \
	--network $(NETWORK_NAME) \
	-e DB_PASSWORD=$(MSSQL_PASS) \
	-v $(PWD)/test/mysql:/tmp/testmysql/:ro \
	$(SYNCDBDOCS_IMAGE) \
	-h $(MSSQL_CONTAINER) \
	-p $(MSSQL_PORT) \
	-u $(MSSQL_USER) \
	-d $(DB_NAME)

test-mssql:
	$(MSSQL_RUN_SYNCDBDOCS) -format=text > /tmp/dbtest.result.txt
	diff $(PWD)/test/mssql/dbtest-from-scratch.expected.txt /tmp/dbtest.result.txt || (echo "MSSQL Test001.txt failed" && false)

SQLITE_RUN_SYNCDBDOCS = docker run --rm \
	--network $(NETWORK_NAME) \
	-e DB_PASSWORD=$(MSSQL_PASS) \
	-v $(PWD)/test/sqlite:/tmp/testsqlite/:ro \
	$(SYNCDBDOCS_IMAGE) \
	-h /tmp/testsqlite/$(SQLITE_FILE) \
	-d $(DB_NAME)

test-sqlite:
	$(SQLITE_RUN_SYNCDBDOCS) -format=text > /tmp/dbtest.result.txt
	diff $(PWD)/test/sqlite/dbtest-from-scratch.expected.txt /tmp/dbtest.result.txt || (echo "SQLITE Test001.txt failed" && false)

