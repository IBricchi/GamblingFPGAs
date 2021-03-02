#!/bin/bash

set -eou pipefail

CUSTOMPORT="${1:-3000}"
CUSTOMDBNAME="${2:-serverDB.db}"

make all

./bin/server -httpPort ${CUSTOMPORT} -db ${CUSTOMDBNAME}
