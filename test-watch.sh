#!/usr/bin/env bash

while true
do
    inotifywait -qq -r -e create,close_write,modify,move,delete --exclude '\.git' ./ \
        && go test ./... \
        ; echo '-----------------------------------------' \
        && date +%H:%M:%S \
        && echo '-----------------------------------------'
done
