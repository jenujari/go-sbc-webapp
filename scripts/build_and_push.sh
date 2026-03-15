#!/usr/bin/env bash
set -euo pipefail

username="jhon5456"
service="swewebapp"

# Get short SHA of current commit
shortSha=$(git rev-parse --short HEAD 2>/dev/null || true)
if [ -z "$shortSha" ]; then
  echo "Failed to get git short SHA. Are you in a git repository?"
  exit 1
fi

sha_image="docker.io/${username}/${service}:${shortSha}"
latest_image="docker.io/${username}/${service}:latest"

# Determine current branch; set MAIN_BRANCH_FLAG=1 when on main
MAIN_BRANCH_FLAG=0
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)
if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
  MAIN_BRANCH_FLAG=1
fi

echo "Determining container CLI to use (prefer 'docker' over 'podman')..."
if command -v docker >/dev/null 2>&1; then
  cli="docker"
elif command -v podman >/dev/null 2>&1; then
  cli="podman"
else
  echo "Neither 'docker' nor 'podman' CLI found in PATH. Install one of them to build images."
  exit 1
fi

echo "Using '${cli}' for image operations."

echo "Building ${sha_image}"
${cli} build -t "${sha_image}" -f Dockerfile .

if [ "$MAIN_BRANCH_FLAG" -eq 1 ]; then
  echo "Tagging ${sha_image} as ${latest_image}"
  ${cli} tag "${sha_image}" "${latest_image}"
else
  echo "Not on main (current: ${current_branch}); skipping tagging as :latest"
fi

if [ -n "${DOCKER_USERNAME:-}" ] && [ -n "${DOCKER_PASSWORD:-}" ]; then
  echo "Logging in to docker.io as ${DOCKER_USERNAME} using ${cli}..."
  echo "${DOCKER_PASSWORD}" | ${cli} login docker.io -u "${DOCKER_USERNAME}" --password-stdin
else
  echo "DOCKER_USERNAME or DOCKER_PASSWORD not set; skipping login (push may fail)."
fi

echo "Pushing ${sha_image} to docker.io"
${cli} push "${sha_image}"

if [ "$MAIN_BRANCH_FLAG" -eq 1 ]; then
  echo "Pushing ${latest_image} to docker.io"
  ${cli} push "${latest_image}"
  echo "Done: ${sha_image} and ${latest_image}"
else
  echo "Done: pushed ${sha_image}; skipped pushing :latest because current branch is ${current_branch}"
fi
