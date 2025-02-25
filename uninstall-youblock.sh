#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo -e "\033[1;31mERROR: This script must be run as root. Use sudo or execute as root user.\033[0m" >&2
    exit 1
fi

# Red YouBlock ASCII Art
echo -e "\033[1;31m"
cat << "EOL"

██╗   ██╗ ██████╗ ██╗   ██╗██████╗ ██╗      ██████╗  ██████╗██╗  ██╗
╚██╗ ██╔╝██╔═══██╗██║   ██║██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝
 ╚████╔╝ ██║   ██║██║   ██║██████╔╝██║     ██║   ██║██║     █████╔╝ 
  ╚██╔╝  ██║   ██║██║   ██║██╔══██╗██║     ██║   ██║██║     ██╔═██╗ 
   ██║   ╚██████╔╝╚██████╔╝██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗
   ╚═╝    ╚═════╝  ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝

EOL
echo -e "\033[0m"
echo -e "\033[1;31m=== YouBlock Service Removal ===\033[0m"
echo

# Stop and remove service
echo -e "\033[1;33m[1/3] Stopping service...\033[0m"
systemctl stop youtube-blocker || true

echo -e "\033[1;33m[2/3] Removing components...\033[0m"
systemctl disable youtube-blocker || true
rm -f /etc/systemd/system/youtube-blocker.service /usr/local/bin/youtube-blocker.sh

echo -e "\033[1;33m[3/3] Cleaning up...\033[0m"
systemctl daemon-reload

# Final message
echo
echo -e "\033[1;31m✔ YouBlock services have been purged\033[0m"
echo -e "\033[1;33mNote: Manual removal of /etc/hosts entries may be required\033[0m"
echo

# Wait for key press
echo -en "\033[1;34mPress any key to conclude uninstallation...\033[0m"
read -n 1 -s -r
echo
