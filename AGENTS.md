# Agent Instructions for go-sbc-webapp

This project is a Go-based web application for financial astrology, serving as a frontend for the `sweapi` gRPC service.

## Tech Stack

- **Backend**: Go (Golang) using standard `net/http` for routing (`http.ServeMux`).
- **Frontend**: Go `html/template`, Vanilla CSS ,Vanilla JS with HTMX, Tailwind CSS and Alpine JS imported using CDN links.
- **Service Communication**: gRPC via `google.golang.org/grpc` (interacts with `sweapi`).
- **Configuration**: Managed by `viper` (`config/` directory).
- **Environment**: Containerized development using `podman compose`.
- **Live Reloading**: Orchestrated via `air` (`go tool air`).
- **Context Management**: Uses `github.com/jenujari/runtime-context` for process lifecycle and logging.

## Dev Environment Tips

- **Developer Workflow**: Run `make dev` to start the infrastructure (`sweapi`) and launch the web app with live reloading.
- **Tooling**: Requires `podman` and `podman-compose`.
- **Interactive Testing**: Use `make grpc-ui` to test the underlying gRPC services.
- **Infrastructure**: The app relies on `sweapi` running via `podman compose`. Ensure the `sweapi` service is healthy if data fetching fails.

## Directory Structure

- `/server`: HTTP handlers and server lifecycle logic.
- `/html/template`: Server-side HTML templates.
- `/html/static`: Client-side assets (CSS/JS).
- `/config`: Configuration schema and initialization.
- `/lib`: Internal library code, including gRPC client wrappers.

## Testing Instructions

- _Note: Currently, there are no unit tests in this repository._
- To verify changes, run `make dev` and manually check the UI or use `curl` to hit the endpoints (`/`, `/pos-table`, etc.).

## PR Instructions

- **Code Style**: Run `go fmt` before committing.
- **Dependencies**: Ensure `go.mod` and `go.sum` are updated if new packages are added.
- **Infra Consistency**: Verify that `compose.yaml` remains compatible with `podman compose`.
- **Template Safety**: Ensure HTML templates are properly escaped and handled.
