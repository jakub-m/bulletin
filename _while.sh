#!/bin/bash

set -eu

while true; do
  date
  if [ "$(date +%u)" -eq 5 ]; then
    git pull
    make up-push
  fi
  echo Zzzzz....
  sleep 14400 # 4 hours
done
