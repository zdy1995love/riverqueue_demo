#!/bin/bash

# PostgreSQL Docker Setup for River Queue Demo
# This script starts a PostgreSQL container with the riverqueue schema

echo "🐳 Starting PostgreSQL container..."

# Configuration
POSTGRES_USER="riverqueue"
POSTGRES_PASSWORD="riverqueue_password"
POSTGRES_DB="river_demo"
POSTGRES_PORT="5432"
CONTAINER_NAME="riverqueue-postgres"

# Start PostgreSQL container
docker run --name ${CONTAINER_NAME} \
  -e POSTGRES_USER=${POSTGRES_USER} \
  -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
  -e POSTGRES_DB=${POSTGRES_DB} \
  -p ${POSTGRES_PORT}:5432 \
  -v $(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql \
  -d postgres:15-alpine

echo "✅ PostgreSQL container started"
echo ""
echo "Connection details:"
echo "  Host: localhost"
echo "  Port: ${POSTGRES_PORT}"
echo "  Database: ${POSTGRES_DB}"
echo "  User: ${POSTGRES_USER}"
echo "  Password: ${POSTGRES_PASSWORD}"
echo ""
echo "📝 Update your setting/config_DEV.jsonc with:"
echo "  \"river_database_url\": \"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?search_path=riverqueue\""
echo ""
echo "🔍 To check container status: docker ps"
echo "🛑 To stop container: docker stop ${CONTAINER_NAME}"
echo "🗑️  To remove container: docker rm ${CONTAINER_NAME}"
