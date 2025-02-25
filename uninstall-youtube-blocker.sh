#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. Use sudo or execute as root user." >&2
    exit 1
fi

systemctl stop youtube-blocker
systemctl disable youtube-blocker
rm /etc/systemd/system/youtube-blocker.service /usr/local/bin/youtube-blocker.sh
systemctl daemon-reload
