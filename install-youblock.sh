#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo -e "\033[1;31mERROR: Root privileges required. Use: sudo ./install-youblock\033[0m" >&2
    exit 1
fi

# YouBlock ASCII Art
echo -e "\033[1;32m"
cat << "EOL"

██╗   ██╗ ██████╗ ██╗   ██╗██████╗ ██╗      ██████╗  ██████╗██╗  ██╗
╚██╗ ██╔╝██╔═══██╗██║   ██║██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝
 ╚████╔╝ ██║   ██║██║   ██║██████╔╝██║     ██║   ██║██║     █████╔╝ 
  ╚██╔╝  ██║   ██║██║   ██║██╔══██╗██║     ██║   ██║██║     ██╔═██╗ 
   ██║   ╚██████╔╝╚██████╔╝██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗
   ╚═╝    ╚═════╝  ╚═════╝ ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝

EOL
echo -e "\033[0m"
echo -e "\033[1;34m=== Website Blocker Service Installation ===\033[0m"
echo

# Create blocker script
echo -e "\033[1;33m[1/4] Deploying blocker engine...\033[0m"
cat > /usr/local/bin/youblock.sh <<'EOL'
#!/bin/bash
while true; do
    for host in youtube.com www.youtube.com; do
        if ! grep -qFx "127.0.0.2 $host" /etc/hosts; then
            echo "127.0.0.2 $host" >> /etc/hosts
        fi
    done
    sleep 5
done
EOL

# Set permissions
echo -e "\033[1;33m[2/4] Configuring security...\033[0m"
chmod 744 /usr/local/bin/youblock.sh

# Create systemd service
echo -e "\033[1;33m[3/4] Building service layer...\033[0m"
cat > /etc/systemd/system/youblock.service <<EOL
[Unit]
Description=YouBlock Website Restriction Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/youblock.sh
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOL

# Enable service
echo -e "\033[1;33m[4/4] Activating protection...\033[0m"
systemctl daemon-reload
systemctl enable youblock
systemctl start youblock

# Completion message
echo
echo -e "\033[1;32m✔ Protection Matrix Activated\033[0m"
echo -e "\033[1;36mContinuous monitoring engaged for restricted domains\033[0m"
echo
echo -e "\033[35mSystem Controls:"
echo -e "  Status Check: systemctl status youblock"
echo -e "  View Logs:    journalctl -u youblock.service"
echo -e "  Integrity Test: Remove any rule from /etc/hosts - auto-restore in 5s\033[0m"

# Wait for key press
echo -en "\033[1;34mPress any key to conclude uninstallation...\033[0m"
read -n 1 -s -r
echo
