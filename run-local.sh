#!/bin/bash

if ! $(docker network ls | grep -q apetito-imageprocessor);
then
  docker network create \
    --driver=bridge \
    --subnet=171.20.0.0/16 \
    --gateway=171.20.0.1 \
    apetito-imageprocessor
fi

exec docker-compose -p "apetito-imageprocessor" up --build -d