#!/bin/sh
set -e

if [ ! -d /var/lib/prometheus-scw-sd ] && [ ! -d /etc/prometheus-scw-sd ]; then
    userdel prometheus-scw-sd 2>/dev/null || true
    groupdel prometheus-scw-sd 2>/dev/null || true
fi
