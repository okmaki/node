#!/bin/bash

CONTAINER=$(cat ./config/container)
MIGRATIONS_DIR=$CONTAINER:/

docker cp ./migrations $MIGRATIONS_DIR
docker exec $CONTAINER /migrations/migrate.sh $1
