#!/usr/bin/env bash
docker rm -f pharmeum-users-api
docker rm -f pharmeum-users-api-migrator
docker rm -f pharmeum-users-postgres

docker rmi -f pharmeumusersapi_pharmeum-users-api-migrator
docker rmi -f pharmeumusersapi_pharmeum-users-api