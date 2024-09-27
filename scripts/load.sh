#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$SCRIPT_DIR/..    
RESULTS_DIR=$ROOT/out/tests/load/results

$SCRIPT_DIR/spinup_bg.sh 

locust -f $ROOT/tests/load/analyze_and_query.py --headless --host http://127.0.0.1:8100 -u 200 -r 20 --run-time 10 --csv $RESULTS_DIR   

$SCRIPT_DIR/teardown.sh 