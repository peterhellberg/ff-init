#!/usr/bin/env bash

killall -q firefly-emulator

firefly_cli build --no-tip && \
firefly_cli emulator -- --id "ff-author-id.ff-app-id"
