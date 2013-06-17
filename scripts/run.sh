#!/bin/bash
cd "$(dirname "$0")/../"
while true
do
    go run *.go -config etc/ircbot.yaml
    sleep 1
done
