#!/bin/bash

fatal() {
    echo $1
    exit 1
}

CONTAINER=$(cat ./config/container)
DATA_VOLUME=$(cat ./config/data_volume)
DATA_MNT=$(cat ./config/data_mnt)

# clear log file

LOG=log.txt
echo "" > $LOG

# create volume for the database's data

if ! docker volume ls | grep -q $DATA_VOLUME;
then
    echo "creating ${DATA_VOLUME}..."

    docker volume create $DATA_VOLUME >> $LOG 2>&1 || fatal "failed to create ${DATA_VOLUME}"
fi

# create the database container

if ! docker container ls | grep -q $CONTAINER;
then
    echo "creating ${CONTAINER}..."

    docker run \
    --name $CONTAINER \
    -v $DATA_MNT:/var/lib/cassandra \
    -p 9042:9042 \
    -d cassandra:latest \
    >> $LOG 2>&1 \
    || fatal "failed to create ${CONTAINER}"
fi

ATTEMPT=0
SUCCESS=false

echo "initializing ${CONTAINER}..."

while [ $ATTEMPT -lt 5 ] && [ $SUCCESS == false ]
do
    if [ $ATTEMPT -ne 0 ]
    then
        echo "waiting 10 seconds before next attempt..."
        sleep 10;
    fi

    ./update.sh --init >> $LOG 2>&1 && SUCCESS=true
    ((ATTEMPT++))

    if [ $SUCCESS == false ]
    then
        echo "attempt ${ATTEMPT} failed"
        continue
    fi
done

if [ $SUCCESS == true ]
then
    echo "DONE"
else
    echo "FAILED"
fi
