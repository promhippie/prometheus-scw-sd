#!/bin/sh
set -e

chown -R prometheus-scw-sd:prometheus-scw-sd /etc/prometheus-scw-sd
chown -R prometheus-scw-sd:prometheus-scw-sd /var/lib/prometheus-scw-sd
chmod 750 /var/lib/prometheus-scw-sd

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload

    if systemctl is-enabled --quiet prometheus-scw-sd.service; then
        systemctl restart prometheus-scw-sd.service
    fi
fi
