#!/bin/bash

set -eou pipefail

CUSTOMPORT="${1:-3000}"
CUSTOMDBNAME="${2:-serverDB.db}"

mkdir -p bin db

make all

./bin/server -httpPort ${CUSTOMPORT} -db ${CUSTOMDBNAME}
