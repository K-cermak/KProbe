#!/bin/sh
set -e

# Stop the service only on removal, not on upgrade (deb: "remove", rpm: 0)
if [ "$1" = "remove" ] || [ "$1" = "0" ]; then
  systemctl stop kprobe.service || true
  systemctl disable kprobe.service || true
fi
