#!/bin/bash

#server
export SRV_PORT=5000

#database connection
export DB_HOST=localhost
export DB_USER=postgres
export DB_NAME=postgres
export DB_PORT=5432
export DB_PASSWORD=postgres
export DB_SSLMODE=disable

#database migration
export M_PATH=file://_migrations
