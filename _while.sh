#!/bin/bash

set -eu

while true; do
  date
  if [ "$(date +%u)" -eq 5 ]; then
    make up-push
  fi
  echo Zzzzz....
  sleep 21600 # 6 hours
done
