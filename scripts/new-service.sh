#!/bin/bash

if [ $# -eq 0 ]
then 
    echo "A name for the new service must be provided!"
    exit 1
fi

SERVICE_NAME=$(echo $1 | grep -E '^[a-z]{3,16}$')

if [ -z $SERVICE_NAME ]
then
    echo "Invalid service name - must be 3 to 16 lowercase letters!"
    exit 1  
fi

SERVICE_PATH=../$SERVICE_NAME

if [ -d $SERVICE_PATH ]
then
    echo "Invalid service name - a service with this name already exists!"
    exit 1  
fi

cp -r ../templates/service $SERVICE_PATH
(cd $SERVICE_PATH && go mod init github.com/okmaki/node/$SERVICE_NAME) || (rm -rf $SERVICE_PATH; echo "Failed to init go module")

echo "Done!"
exit 0
