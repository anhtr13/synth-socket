#!/bin/bash

sqlc generate

targets=("./api/internal/" "./socket/internal/" "./cron/internal/")

for target in ${targets[@]}; do
  mkdir -p $target
done

for target in ${targets[@]}; do
  cp -r -f ./share/* $target
done
