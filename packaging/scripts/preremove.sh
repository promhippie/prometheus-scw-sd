#!/bin/sh
set -e

systemctl stop prometheus-scw-sd.service || true
systemctl disable prometheus-scw-sd.service || true
