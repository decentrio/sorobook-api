#!/usr/bin/env bash

# --------------
# Commands to run locally
#

set -eo pipefail

echo "Generating proto code"

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    # Check if the go_package in the file is pointing to realionetwork
    if grep -q "option go_package.*decentrio" "$file"; then
      buf generate --template proto/buf.gen.yaml "$file"
    fi
  done
done
