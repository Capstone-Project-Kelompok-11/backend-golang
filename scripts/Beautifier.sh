#!/usr/bin/env bash

#set -xe

# shellcheck disable=SC2046
# shellcheck disable=SC2002
find app -name '*.go' -type f -print0 | while IFS= read -r -d '' file
do

  echo "$file"
  cat<<<$(cat "$file" | sed -e 's/\r\n/\n/g' -e 's/\t/  /g' -e 's/\n\n/\n/g')>$file
done