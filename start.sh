#!/usr/bin/env bash
set -e

status() { echo -e "\033[1;36m$1\033[0m"; }
error()  { echo -e "\033[1;31m$1\033[0m"; }

# Ensure Docker daemon is up
if ! command -v docker >/dev/null 2>&1; then
  error "Docker CLI not found. Install Docker Desktop and retry."
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  error "Docker daemon is not running. Start Docker Desktop and retry."
  exit 1
fi

status "Running: docker compose up -d --build"
docker compose up -d --build

echo ""
echo "🚀 Containers started successfully!"
echo "-----------------------------------"
echo "🌐 Frontend:   http://localhost:3000"
echo "🛠 Backend:    http://localhost:8081"
echo "📊 Hasura:     http://localhost:8082"
echo "🐘 Postgres:   host 5433 → container 5432"
echo "-----------------------------------"

open_url() {
  local url="$1"
  if command -v xdg-open >/dev/null 2>&1; then xdg-open "$url" >/dev/null 2>&1 || true; fi
  if command -v open >/dev/null 2>&1; then open "$url" >/dev/null 2>&1 || true; fi
  # Git Bash / MSYS often supports 'start'
  if command -v start >/dev/null 2>&1; then start "" "$url" >/dev/null 2>&1 || true; fi
}

# Try to open the endpoints
open_url "http://localhost:3000"
open_url "http://localhost:8082"
open_url "http://localhost:8081/health"
