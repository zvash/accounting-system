#!/bin/sh

set -e

echo "migrate db"
/app/migrate -path sql/migration -database "$DB_SOURCE" -verbose up
echo "start the app"
exec "$@"