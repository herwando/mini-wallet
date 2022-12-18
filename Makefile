MODULES      = $(shell cd module && ls -d */)
TMP_DIR     := $(shell mktemp -d)
UNAME       := $(shell uname)

# Default postgres migration settings
export POSTGRES_USER ?= root
export POSTGRES_PASS ?= rootpw
export POSTGRES_HOST ?= localhost
export POSTGRES_PORT ?= 15432
export POSTGRES_DATABASE ?= wallets
export POSTGRES_QUERYSTRING ?= sslmode=disable
export POSTGRES_MIGRATEUP ?= up
export POSTGRES_MIGRATEDOWN ?= down

dep:
	GO111MODULE=on go mod download
	GO111MODULE=on go mod verify
	GO111MODULE=on go mod tidy

run:
	go build -o api cmd/*.go
	./api

docker-up:
	@docker-compose -f dev/docker-compose.yml up -d

docker-down:
	@docker-compose -f dev/docker-compose.yml down

bin:
	@mkdir -p bin

tool-migrate: bin
ifeq ($(UNAME), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(UNAME), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo "Your OS is not supported."
endif

migrate-up:
	@$(foreach module, $(MODULES), cp module/$(module)/db/migrate/*.sql $(TMP_DIR);)
	@bin/migrate -source file://$(TMP_DIR) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASS)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?$(POSTGRES_QUERYSTRING)" $(POSTGRES_MIGRATEUP)

migrate-down:
	@$(foreach module, $(MODULES), cp module/$(module)/db/migrate/*.sql $(TMP_DIR);)
	@bin/migrate -source file://$(TMP_DIR) -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASS)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?$(POSTGRES_QUERYSTRING)" $(POSTGRES_MIGRATEDOWN)
