#!/bin/bash

if [ -z "$KAFKA_DIR" ]
then
      echo "\$KAFKA_DIR is empty"
      exit 1
else
      echo "Going to start kafka from $KAFKA_DIR"
fi

go get github.com/mattn/goreman

cp kafka_server_jaas.conf $KAFKA_DIR
cp server.1.properties $KAFKA_DIR
cp server.2.properties $KAFKA_DIR
cp server.3.properties $KAFKA_DIR
cp Procfile $KAFKA_DIR
cd $KAFKA_DIR && goreman start