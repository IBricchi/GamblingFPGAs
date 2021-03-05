#!/bin/bash

set -eou pipefail

CUSTOMDBNAME="${1:-serverDB.db}"

mkdir -p bin db

make credentials

./bin/credentials -db ${CUSTOMDBNAME}
