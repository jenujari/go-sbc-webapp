.PHONY: build-and-push-container run-infra stop-infra dev build run-dev run-prod grpc-ui

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

POSTGRESQL_URL := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)&search_path=public"

run-infra:
	podman-compose up -d sweapi

stop-infra:
	podman-compose down -v

dev: run-infra
	go tool air

build:
	podman-compose build

run-dev:
	podman-compose up --build webapp-dev

run-prod:
	podman-compose up -d webapp

grpc-ui:
	podman run --rm --network=host -p 8080:8080 docker.io/fullstorydev/grpcui -plaintext localhost:5678

t1:
	@echo $(POSTGRESQL_URL)

# make migrate-create name=test_migration
migrate-create:
	go tool migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate -database $(POSTGRESQL_URL) -path db/migrations up

migrate-down:
	go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate -database $(POSTGRESQL_URL) -path db/migrations down $(N)

build-and-push-container:
	./scripts/build_and_push.sh

jenkins-deploy:
	./scripts/trigger_deploy.sh
