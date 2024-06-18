#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen

cd proto
echo "Generate sorobook swagger files"
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    echo "$query_file"
    buf generate --template buf.gen.swagger.yaml "$query_file"
  fi
done
rm -rf github.com

cd ..
echo "Combine swagger files"
# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./docs/swagger-ui/config.json -o ./docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen

echo "Update statik data"
install_statik() {
  go install github.com/rakyll/statik@v0.1.7
}
install_statik

# generate binary for static server
statik -f -src=./docs/swagger-ui -dest=./docs