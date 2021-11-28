#!/bin/bash

fatal() {
    echo $1
    exit 1
}

if [ -z $1 ]
then
    fatal "timestamp is required"
fi

TIMESTAMP=$1

if ! echo "$TIMESTAMP" | egrep -qo "[0-9]{12}";
then
    fatal "timestamp is not in a valid format yyyyMMddHHmm"
fi

SCRIPT=$2

if [ -z $SCRIPT ]
then
    fatal "script name is required"
fi

echo "\
INSERT INTO db.migrations (cluster, timestamp, script)\
VALUES ('log', '${TIMESTAMP}', '${SCRIPT}');" > record_migration.cql

cqlsh -f record_migration.cql || fatal "failed to record migraton"

exit 0
