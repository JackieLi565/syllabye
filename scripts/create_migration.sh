#!/bin/bash

# create_migration.sh creates a migration in the migrations directory.

if [ $# -eq 0 ]; then
  echo "Error: No migration name provided."
  echo "Usage: $0 <migration_name>"
  exit 1
fi

migrate create -ext sql -dir migrations -seq $1
