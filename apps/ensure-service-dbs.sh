#!/bin/sh
set -e

echo "Ensuring device_db and telemetry_db exist..."
until pg_isready -h postgres -U postgres >/dev/null 2>&1; do
  sleep 1
done

psql -h postgres -U postgres -d postgres -v ON_ERROR_STOP=1 <<'SQL'
SELECT 'CREATE DATABASE device_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'device_db')\gexec

SELECT 'CREATE DATABASE telemetry_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'telemetry_db')\gexec
SQL

echo "Service databases are ready."
