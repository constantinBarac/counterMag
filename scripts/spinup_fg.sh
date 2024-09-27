#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$SCRIPT_DIR/..
EXECUTABLE_PATH=$ROOT/bin/countermag

run_master () {
    $EXECUTABLE_PATH --cluster 127.0.0.1:8000 --port 8000 > /dev/null
}

run_slave () {
    $EXECUTABLE_PATH --cluster 127.0.0.1:8000 --port 8001 > /dev/null
}


echo "Starting master in background..."
run_master &
sleep 2
echo "Starting slave..."
run_slave
