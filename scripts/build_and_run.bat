@echo off
REM This script will:
REM 1. Build all images defined in docker-compose.yml
REM 2. Bring all services up in detached mode
REM 3. Apply database migrations via the rating-service-migrator container

echo Step 1) Stopping and removing existing containers (if any)...
docker compose down -v

echo Step 2) Building all images...
docker compose build

echo Step 3) Starting all services in detached mode...
docker compose up -d

REM Optional: Wait a few seconds for services like PostgreSQL to be ready.
echo Step 4) Waiting for services to initialize...
timeout /t 5 >nul

REM Run migrations
echo Step 5) Running database migrations...
docker compose run --rm rating-service-migrator

echo All done!
echo Use "docker compose logs -f" to follow service logs.
echo Available services:
echo  - Rating Service:       http://localhost:8181/swagger/index.html
echo  - Notification Service: http://localhost:9191/swagger/index.html
echo  - Kafka UI:            http://localhost:8082
