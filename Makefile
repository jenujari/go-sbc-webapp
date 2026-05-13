.PHONY: build-and-push-container run-infra stop-infra dev build run-dev run-prod grpc-ui up-infra up-dev up-prod down

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

COMPOSE := podman compose
PROJECT_NAME ?= $(notdir $(CURDIR))
DEFAULT_NETWORK := $(PROJECT_NAME)_default
POSTGRESQL_URL := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)&search_path=public"

t1:
	@echo $(POSTGRESQL_URL)

build:
	$(COMPOSE) --profile build build

up-infra:
	@if podman network exists $(DEFAULT_NETWORK); then \
		network_label=$$(podman network inspect $(DEFAULT_NETWORK) --format '{{ index .Labels "com.docker.compose.network" }}' 2>/dev/null || true); \
		if [ "$$network_label" != "default" ]; then \
			echo "Removing stale compose network '$(DEFAULT_NETWORK)' (label: '$$network_label')"; \
			podman network rm $(DEFAULT_NETWORK); \
		fi; \
	fi
	$(COMPOSE) --profile infra up -d

up-dev:
	$(COMPOSE) --profile dev up --build

up-prod:
	$(COMPOSE) --profile prod up -d --build

down:
	$(COMPOSE) down -v

dev: up-infra
	go tool air

grpc-ui:
	podman run --rm --network=host -p 8080:8080 docker.io/fullstorydev/grpcui -plaintext localhost:5678

# make migrate-create name=test_migration
migrate-create:
	go tool migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate -database $(POSTGRESQL_URL) -path db/migrations up

migrate-down:
	go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate -database $(POSTGRESQL_URL) -path db/migrations down $(N)
