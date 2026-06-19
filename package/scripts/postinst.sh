#!/bin/sh
set -e

# Initialize DB if it doesn't exist yet (first install)
if [ ! -f /opt/kprobe/db.sqlite ]; then
  echo "y" | /opt/kprobe/cli db init
fi

systemctl daemon-reload
systemctl enable kprobe.service
systemctl start kprobe.service
sysctl --system
