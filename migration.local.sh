#!/bin/bash

export MIGRATION_DIR=./migrations
export DB_HOST="localhost"
export DB_PORT="54321"
export DB_NAME="url_shortener"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_SSL=disable

export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSL}"

goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" up -v