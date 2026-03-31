#!/usr/bin/env bash
set -euo pipefail

# List available iOS Simulator destinations
# Usage: list_destinations.sh

xcrun simctl list devices available | grep -E "iPhone|iPad" || echo "No simulators found. Run 'xcrun simctl list devices' for more details."

