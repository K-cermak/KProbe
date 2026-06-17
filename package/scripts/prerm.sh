#!/bin/sh
set -e
systemctl stop kprobe.service || true
systemctl disable kprobe.service || true
