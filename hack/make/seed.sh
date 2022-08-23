#!/usr/bin/env bash
set -e

echo 'Waiting for the database to become responsive...'
./hack/wait-for-it.sh localhost:5432 -q -t 30

echo 'Seeding the database with sample data...'
psql -U postgres -c 'DROP DATABASE IF EXISTS homework'
psql -U postgres <./data/cpu_usage.sql
psql -U postgres -d homework -c '\COPY cpu_usage FROM ./data/cpu_usage.csv CSV HEADER'

echo 'Done.'
