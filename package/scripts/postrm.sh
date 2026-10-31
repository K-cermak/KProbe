#!/bin/sh
set -e
# deb: "remove", rpm: 0
if [ "$1" = "remove" ] || [ "$1" = "0" ]; then
  systemctl daemon-reload
  sysctl --system
fi
