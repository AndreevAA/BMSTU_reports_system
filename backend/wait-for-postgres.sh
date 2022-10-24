#!/bin/sh
# wait-for-postgres.sh

set -e
  
host="$1"
shift
  
until PGPASSWORD="reports_system" psql -h "$host" -U "reports_system" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done
  
>&2 echo "Postgres is up - executing command"
exec "$@"