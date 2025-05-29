#!/bin/bash

set -eux

while true; do
  date
  if [ "$(date +%u)" -eq 5 ]; then
    make up-push
  fi
  sleep 21600 # 6 hours
done
