#!/usr/bin/env bash
set -euo pipefail

# Determine current branch; set MAIN_BRANCH_FLAG=1 when on main
MAIN_BRANCH_FLAG=0
current_branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)
if [ "$current_branch" = "main" ]; then
  MAIN_BRANCH_FLAG=1
fi

if [ "$MAIN_BRANCH_FLAG" -eq 0 ]; then
  echo "Not on main branch (current: ${current_branch}); skipping Jenkins trigger."
  exit 0
fi

# Get short SHA of current commit
shortSha=$(git rev-parse --short HEAD 2>/dev/null || true)
if [ -z "$shortSha" ]; then
  echo "Failed to get git short SHA. Are you in a git repository?"
  exit 1
fi

if [ -n "${JENKINS_TOKEN:-}" ]; then
  echo "Triggering Jenkins build with SHA ${shortSha:-}"
  curl -X POST -u github-trigger:${JENKINS_TOKEN} "https://jenkins.vilabs.co.in/job/sbc-web-app-deployment/buildWithParameters?token=trigger-token-sbc-webapp&IMAGE_TAG=latest&IMAGE_SHA=${shortSha}"
else
  echo "JENKINS_TOKEN not set; skipping Jenkins trigger."
fi