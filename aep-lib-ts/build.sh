#!/usr/bin/env bash
set -euo pipefail

TARGET_DIR="${BUILD_WORKSPACE_DIRECTORY:-$(cd "$(dirname "$0")" && pwd)}"
cd "${TARGET_DIR}"

if [ ! -d "node_modules" ]; then
  npm ci --silent
fi

npx tsc
