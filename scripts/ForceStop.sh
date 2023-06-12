#!/usr/bin/env bash

PID=$(ps -o args,pid -C start | pgrep -i start | rev | awk '{print $1}' | rev)

if [ -n "$PID" ]; then

  kill -9 "$PID"
  echo Stopped ...
fi