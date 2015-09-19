#!/bin/bash

echo "Waiting for Consul join..."

TRIES=10

until consul join $1; do
  let TRIES-=1

  echo "Consul join failed. Retries left: $TRIES"
  if [ $TRIES -le 0 ]; then
    # Update the "magic" key in Consul that says "no deploys!"
    echo "Could not reach application, assuming it failed to start!"
    exit 27
  fi

  sleep 2
done
