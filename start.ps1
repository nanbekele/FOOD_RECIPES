# Start all services and open endpoints
param(
  [switch]$NoBrowser
)

Write-Host "Running: docker compose up -d --build" -ForegroundColor Cyan
try {
  docker version | Out-Null
} catch {
  Write-Error "Docker Desktop is not running. Please start Docker Desktop and retry."
  exit 1
}

$proc = Start-Process -FilePath "docker" -ArgumentList "compose up -d --build" -NoNewWindow -PassThru -Wait
if ($proc.ExitCode -ne 0) {
  Write-Error "docker compose failed with exit code $($proc.ExitCode)"
  exit $proc.ExitCode
}

Write-Host "" 
Write-Host "Containers started successfully!" -ForegroundColor Green
Write-Host "-----------------------------------"
Write-Host "Frontend:   http://localhost:3000"
Write-Host "Backend:    http://localhost:8081"
Write-Host "Hasura:     http://localhost:8082"
Write-Host "Postgres:   host 5433 -> container 5432"
Write-Host "-----------------------------------"

if (-not $NoBrowser) {
  Start-Process "http://localhost:3000" | Out-Null
  Start-Process "http://localhost:8082" | Out-Null
  Start-Process "http://localhost:8081/health" | Out-Null
}
