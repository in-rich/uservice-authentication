#!/bin/bash

int_handler()
{
    docker compose -f docker-compose.test.yml down
}
trap int_handler INT

docker compose -f docker-compose.test.yml up -d

firebase emulators:exec "go run gotest.tools/gotestsum@latest --format pkgname -- -p 1 ./..."

docker compose -f docker-compose.test.yml down
