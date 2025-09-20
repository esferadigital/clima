#!/usr/bin/env bash

watchexec \
    --restart \
    --stop-timeout 3s \
    --wrap-process session \
    --clear \
    --exts go,mod,sum \
    -- "go run main.go --debug"

