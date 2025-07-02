#!/bin/sh
set -e

if ! getent group prometheus-scw-sd >/dev/null 2>&1; then
    groupadd --system prometheus-scw-sd
fi

if ! getent passwd prometheus-scw-sd >/dev/null 2>&1; then
    useradd --system --create-home --home-dir /var/lib/prometheus-scw-sd --shell /bin/bash -g prometheus-scw-sd prometheus-scw-sd
fi
