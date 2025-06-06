#!/usr/bin/env bash
set -e

export PGDATA=/var/lib/postgresql/data
export PORT=${PORT:-8080}
export FERN_REPORTER_DB_DSN=${FERN_REPORTER_DB_DSN:-"postgres://fern:fern@localhost:5432/fern?sslmode=disable"}

echo "=== init-postgres.sh invoked ==="

# 1) Initialize DB on first run
if [ ! -f "$PGDATA/PG_VERSION" ]; then
  echo "Initializing PostgreSQL database..."
  initdb -D "$PGDATA"
fi

# 2) Start Postgres in background
echo "Starting PostgreSQL..."
pg_ctl -D "$PGDATA" -l "$PGDATA/logfile" start

# 3) Wait for server to be ready
echo "Waiting for PostgreSQL to come online..."
until pg_isready -q; do
  sleep 1
done

# 4) Create 'fern' user & DB if missing
echo "Creating fern user & database if needed..."
psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
  DO \$\$
  BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'fern') THEN
      CREATE ROLE fern WITH LOGIN PASSWORD 'fern';
    END IF;
  END
  \$\$;
  
  SELECT 'CREATE DATABASE fern OWNER fern'
  WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'fern')\gexec
EOSQL

# 5) Apply migrations
echo "Applying migrations..."
migrate -path /workspace/../fern-reporter/migrations -database "$FERN_REPORTER_DB_DSN" up || true

# 6) Launch MCP server (PID 1)
echo "Starting MCP server on port $PORT..."
exec go run cmd/serve.go --port="$PORT"
