#!/usr/bin/env bash

if [ -n "$1" ]; then

  cat $1 | grep -i 'func (' | cut -d\) -f2- | cut -d\{ -f1
fi