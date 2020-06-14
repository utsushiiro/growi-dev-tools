#!/bin/bash

if (( $# != 1 )); then
  echo "Error: specify only one argument" 1>&2
  exit 1
fi

GROWI_PID=$1

kill -USR2 $GROWI_PID

./stress-test --rate 50 --duration 610 --random-page-access > ramdom-page-access.log & 
./stress-test --rate 10  --duration 610 --random-page-access > ramdom-page-update.log &

sleep 600

kill -USR2 $GROWI_PID
