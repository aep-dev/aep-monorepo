#!/usr/bin/env python3
"""Manage and verify internal monorepo dependencies for aep-lib-go consumers."""

import argparse
import os
import re
import subprocess
import sys
from pathlib import Path

DEPENDENT_MODULES = ["aepcli", "aepc", "aep-e2e-validator", "terraform-provider-aep"]
ROOT_DIR = Path(__file__).resolve().parent.parent


def run_cmd(cmd: list[str], check: bool = True) -> str:
    """Run a shell command and return its stdout stripped."""
    result = subprocess.run(
        cmd,
        cwd=ROOT_DIR,
        capture_output=True,
        text=True,
        check=check,
    )
    return result.stdout.strip()


def check_deps(base_ref: str) -> int:
    """Verify that dependent modules updated go.mod when aep-lib-go changed."""
    print(f"Checking if aep-lib-go changes require dependent go.mod updates against '{base_ref}'...")

    # Ensure base reference exists
    try:
        run_cmd(["git", "rev-parse", "--verify", base_ref])
    except subprocess.CalledProcessError:
        print(f"Warning: Base ref '{base_ref}' not found locally. Attempting to fetch...")
        subprocess.run(["git", "fetch", "origin", "main"], cwd=ROOT_DIR, capture_output=True)

    # Check for changes in aep-lib-go
    changed_lib_raw = run_cmd(["git", "diff", "--name-only", base_ref, "--", "aep-lib-go/"], check=False)
    changed_lib = [f for f in changed_lib_raw.splitlines() if f.strip()]

    if not changed_lib:
        print(f"No changes in aep-lib-go detected relative to '{base_ref}'. Check passed.")
        return 0

    print("Detected changes in aep-lib-go:")
    for file in changed_lib:
        print(f"  {file}")
    print("\nVerifying that dependent modules have updated their go.mod files...")

    failed = False
    for dep in DEPENDENT_MODULES:
        changed_mod_raw = run_cmd(["git", "diff", "--name-only", base_ref, "--", f"{dep}/go.mod"], check=False)
        changed_mod = [f for f in changed_mod_raw.splitlines() if f.strip()]

        if not changed_mod:
            print(f"Error: '{dep}/go.mod' was not modified alongside aep-lib-go.", file=sys.stderr)
            failed = True
        else:
            print(f"OK: '{dep}/go.mod' was modified.")

    if failed:
        print("\nCheck failed: Changes in aep-lib-go require updating go.mod in all dependent modules.", file=sys.stderr)
        print("Run './scripts/sync-deps.py' (or 'just fix') to automatically synchronize dependencies.", file=sys.stderr)
        return 1

    print("All dependent go.mod files are properly updated.")
    return 0


def sync_deps() -> int:
    """Update require directive in dependent go.mod files to point to current commit."""
    print("Synchronizing internal monorepo dependencies...")

    commit_sha = run_cmd(["git", "rev-parse", "--short=12", "HEAD"])
    commit_time = run_cmd(["git", "log", "-1", "--format=%cd", "--date=format:%Y%m%d%H%M%S", "HEAD"])
    pseudo_version = f"v0.0.0-{commit_time}-{commit_sha}"

    print(f"Target aep-lib-go version: {pseudo_version}")

    pattern = re.compile(r"(github\.com/aep-dev/aep-monorepo/aep-lib-go\s+)v[0-9A-Za-z.\-]+")

    for dep in DEPENDENT_MODULES:
        go_mod_path = ROOT_DIR / dep / "go.mod"
        if not go_mod_path.exists():
            continue

        content = go_mod_path.read_text()
        if "github.com/aep-dev/aep-monorepo/aep-lib-go" in content:
            new_content = pattern.sub(rf"\g<1>{pseudo_version}", content)
            go_mod_path.write_text(new_content)
            print(f"  Updated {dep}/go.mod to {pseudo_version}")
        else:
            print(f"  Skipping {dep} (not directly required in go.mod)")

    print("Synchronization complete. Remember to commit the modified go.mod files.")
    return 0


def main() -> None:
    parser = argparse.ArgumentParser(description="Manage and verify internal monorepo dependencies.")
    parser.add_argument(
        "--check",
        action="store_true",
        help="Verify that dependent go.mod files were updated if aep-lib-go changed.",
    )
    parser.add_argument(
        "--base",
        default=os.environ.get("BASE_REF", "origin/main"),
        help="Specify base git reference for comparison (default: origin/main or $BASE_REF).",
    )

    args = parser.parse_args()

    if args.check:
        sys.exit(check_deps(args.base))
    else:
        sys.exit(sync_deps())


if __name__ == "__main__":
    main()
