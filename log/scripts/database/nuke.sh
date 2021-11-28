#!/bin/bash

CONTAINER=$(cat ./config/container)
DATA_VOLUME=$(cat ./config/data_volume)
DATA_MNT=$(cat ./config/data_mnt)

docker container stop $CONTAINER
docker container rm $CONTAINER
docker volume rm $DATA_VOLUME
rm -rf $DATA_MNT