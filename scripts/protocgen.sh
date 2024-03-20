#!/usr/bin/env bash
cd proto
echo "Formatting protobuf files"

set -e

echo "Generating proto code"
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

for dir in $proto_dirs; do
  buf generate --path="$dir" --template buf.gen.yaml --config buf.yaml
done

go mod tidy
