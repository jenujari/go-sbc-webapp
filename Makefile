.PHONY: build-and-push-container run-infra stop-infra dev build run-dev run-prod grpc-ui

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

build-and-push-container:
	./scripts/build_and_push.sh

jenkins-deploy:
	./scripts/trigger_deploy.sh
