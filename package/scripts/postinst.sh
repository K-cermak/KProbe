#!/bin/sh
set -e
systemctl daemon-reload
systemctl enable kprobe.service
systemctl start kprobe.service
sysctl --system
