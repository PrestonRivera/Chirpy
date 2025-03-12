#!/bin/bash
set -e

echo "Running migrations..."
docker run --rm --network host \
  -v "$(pwd)/sql/schema:/migrations" \
  -e GOOSE_DRIVER="postgres" \
  -e GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgres dbname=chirpy sslmode=disable" \
  ghcr.io/kukymbr/goose-docker:3.19.2 up

echo "Migrations completed successfully"