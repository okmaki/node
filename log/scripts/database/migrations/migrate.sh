#!/bin/bash

fatal() {
    echo $1
    echo $1 >> log.txt
    exit 1
}

if [ "$1" == "--init" ]
then
    echo "running init db script..."
    cqlsh -f /migrations/init_db.cql || fatal "failed to run init db script"
fi

TIMESTAMP_QUERY=$(cqlsh -f /migrations/timestamp.cql)
TIMESTAMP=$(echo "${TIMESTAMP_QUERY}" | egrep -o "[0-9]{12}")
MIGRATIONS=$(ls /migrations)

echo "last migration: ${TIMESTAMP}"

echo "running migration scripts:"

MIGRATION_RUN=false

for SCRIPT in $MIGRATIONS
do
    SCRIPT_TIMESTAMP=$(echo "${SCRIPT}" | egrep -o "[0-9]{12}")

    if [ -z $SCRIPT_TIMESTAMP ]
    then
        continue
    fi

    if [ "$SCRIPT_TIMESTAMP" -gt "$TIMESTAMP" ]
    then
        echo "running ${SCRIPT}"

        cqlsh -f /migrations/$SCRIPT || fatal "failed to run ${SCRIPT}"
        /migrations/record_migration.sh $SCRIPT_TIMESTAMP $SCRIPT || fatal "failed to record ${SCRIPT}"
        
        MIGRATION_RUN=true
    fi
done

if $MIGRATION_RUN;
then
    echo "done"
else
    echo "no migrations were run"
fi
