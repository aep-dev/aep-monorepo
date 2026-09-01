#!/usr/bin/env bash
set -euo pipefail

TARGET_DIR="${BUILD_WORKSPACE_DIRECTORY:-${TEST_SRCDIR}/${TEST_WORKSPACE:-_main}/aep-lib-ts}"
cd "${TARGET_DIR}"

if [ ! -d "node_modules" ]; then
  npm ci --silent
fi

npx jest --rootDir=. --passWithNoTests
