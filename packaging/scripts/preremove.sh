#!/bin/sh
set -e

systemctl stop prometheus-scw-sd.service || true
systemctl disable prometheus-scw-sdpi.service || true
