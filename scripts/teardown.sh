#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$SCRIPT_DIR/.. 

echo "Cleaning up..."

pkill -x countermag 
rm $ROOT/counter-*

echo "Cleanup complete"