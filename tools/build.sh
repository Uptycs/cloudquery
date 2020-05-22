#!/usr/bin/env bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $SCRIPT_DIR/lib.sh

if [[ -z "$SKIP_TESTS" ]]; then RUN_TESTS=true; else RUN_TESTS=false; fi
build $RUN_TESTS
exit 0
