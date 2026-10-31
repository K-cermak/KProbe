#!/bin/sh
set -e

# Initialize DB if it doesn't exist yet (first install)
if [ ! -f /opt/kprobe/db.sqlite ]; then
  echo "y" | /opt/kprobe/cli db init
fi

systemctl daemon-reload
systemctl enable kprobe.service
# Restart (not start) so an upgrade replaces the still running old binary
systemctl restart kprobe.service
sysctl --system
