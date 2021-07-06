#
# Creates some databases, fill them with data, and then run syncdbdocs apps on
# them to check if works as expected
#
RANDVAR      ?= default
NETWORK_NAME ?= sdd_network_$(RANDVAR)

SYNCDBDOCS_IMAGE ?= syncdbdocs

PG_CONTAINER        ?= sdd_pg_$(RANDVAR)
MYSQL_CONTAINER     ?= sdd_msyql_$(RANDVAR)

PG_IMAGE_VERSION    ?= latest
MYSQL_IMAGE_VERSION ?= latest
# MSSQL_IMAGE_VERSION ?= latest

DB_NAME = dbtest
DB_USER = root
DB_PASS = pass

PG_PORT = 5432

.PHONY: all test clean

all:   | up build test clean
clean: | down

up:    | net-up db-up migrate
down:  | db-down net-down

net-up:
	docker network create $(NETWORK_NAME) || true

net-down:
	docker network rm $(NETWORK_NAME) || true

db-up:
	docker run -d --rm \
	  --name $(PG_CONTAINER) \
	  --network $(NETWORK_NAME) \
	  -e POSTGRES_DB=$(DB_NAME) \
	  -e POSTGRES_USER=$(DB_USER) \
	  -e POSTGRES_PASSWORD=$(DB_PASS) \
	  postgres:$(PG_IMAGE_VERSION)

migrate:
	docker run --rm \
	  --network $(NETWORK_NAME) \
		-v $(PWD)/test/postgres:/flyway/sql \
		flyway/flyway:7.11 \
		-url=jdbc:postgresql://$(PG_CONTAINER):$(PG_PORT)/$(DB_NAME) \
		-user=$(DB_USER) \
		-password=$(DB_PASS) \
		-connectRetries=50 \
		migrate

db-down:
	docker stop $(PG_CONTAINER) || true

build:
	docker build . --tag $(SYNCDBDOCS_IMAGE)


RUN_SYNCDBDOCS = docker run --rm \
	--network $(NETWORK_NAME) \
	-e DB_PASSWORD=$(DB_PASS) \
	-v $(PWD)/test/postgres:/tmp/testpg/:ro \
	$(SYNCDBDOCS_IMAGE) \
	-h $(PG_CONTAINER) \
	-p $(PG_PORT) \
	-u $(DB_USER) \
	-d $(DB_NAME)

test:
	$(RUN_SYNCDBDOCS) -format=md > /tmp/dbtest.result.md
	$(RUN_SYNCDBDOCS) -format=text > /tmp/dbtest.result.txt
	diff $(PWD)/test/postgres/dbtest-from-scratch.expected.md /tmp/dbtest.result.md || (echo "Test001.md failed" && false)
	diff $(PWD)/test/postgres/dbtest-from-scratch.expected.txt /tmp/dbtest.result.txt || (echo "Test001.txt failed" && false)

	$(RUN_SYNCDBDOCS) -i /tmp/testpg/dbtest-preserve-order.input > /tmp/dbtest.result
	diff $(PWD)/test/postgres/dbtest-preserve-order.expected.txt /tmp/dbtest.result || (echo "Test002 failed" && false)
