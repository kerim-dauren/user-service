#!/bin/sh

set -e

# Set default migration directory if not provided
: "${MIGRATION_DIR:=db/migration/scripts}"

# Validate required environment variables
if [ -z "$DB_URL" ]; then
  echo "Error: DB_URL is not set. Please set the DB_URL environment variable."
  exit 1
fi

# Check if migration binary exists
if [ ! -f "./db/migration/goose" ]; then
  echo "Warning: Migration tool not found at ./db/migration/goose"
  echo "Migrations will be skipped. Make sure to include the migration tool in your Docker image if needed."
else
  echo "Running database migrations from $MIGRATION_DIR..."

  # Run migrations
  ./db/migration/goose -dir "$MIGRATION_DIR" postgres "$DB_URL" up

  # Check migration result
  if [ $? -ne 0 ]; then
    echo "Error: Database migrations failed."
    exit 1
  fi

  echo "Migrations completed successfully."
fi

# Execute the main application
exec ./main