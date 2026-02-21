

run-infra:
	podman compose up -d sweapi

stop-infra:
	podman compose down

dev: run-infra	
	go tool air

build:
	podman compose build

run: build
	podman compose up

grpc-ui:
	podman run --rm --network=host -p 8080:8080 docker.io/fullstorydev/grpcui -plaintext localhost:5678