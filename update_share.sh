#!/bin/bash

target=("./api/internal/" "./socket/internal/" "./cron/internal/")

sqlc generate

for source in ${target[@]}; do
  mkdir -p $source
  cp -r -f ./share/* $source
done
