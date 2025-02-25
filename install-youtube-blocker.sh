#!/bin/bash
set -e

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. Use sudo or execute as root user." >&2
    exit 1
fi

# Create the blocker script
cat > /usr/local/bin/youtube-blocker.sh <<'EOL'
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

# Set proper permissions for the blocker script
chmod 744 /usr/local/bin/youtube-blocker.sh

# Create systemd service file
cat > /etc/systemd/system/youtube-blocker.service <<EOL
[Unit]
Description=YouTube Blocker Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/youtube-blocker.sh
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOL

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable youtube-blocker
systemctl start youtube-blocker

echo "Installation complete. YouTube blocker daemon is now running."
echo "Check status with: systemctl status youtube-blocker"
echo "Check logs with: journalctl -u youtube-blocker.service"
