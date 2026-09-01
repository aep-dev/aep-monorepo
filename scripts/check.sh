#!/usr/bin/env bash
set -euo pipefail

WORKSPACE="${BUILD_WORKSPACE_DIRECTORY:-$(cd "$(dirname "$0")/.." && pwd)}"
python3 "${WORKSPACE}/scripts/sync-deps.py" --check "$@"

