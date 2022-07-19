PROMETHEUS_SCW_LOG_LEVEL
: Only log messages with given severity, defaults to `info`

PROMETHEUS_SCW_LOG_PRETTY
: Enable pretty messages for logging, defaults to `false`

PROMETHEUS_SCW_WEB_ADDRESS
: Address to bind the metrics server, defaults to `0.0.0.0:9000`

PROMETHEUS_SCW_WEB_PATH
: Path to bind the metrics server, defaults to `/metrics`

PROMETHEUS_SCW_WEB_CONFIG
: Path to web-config file

PROMETHEUS_SCW_OUTPUT_ENGINE
: Enabled engine like file or http, defaults to `file`

PROMETHEUS_SCW_OUTPUT_FILE
: Path to write the file_sd config, defaults to `/etc/prometheus/scw.json`

PROMETHEUS_SCW_OUTPUT_REFRESH
: Discovery refresh interval in seconds, defaults to `30`

PROMETHEUS_SCW_CHECK_INSTANCE
: Enable instance gathering, defaults to `true`

PROMETHEUS_SCW_CHECK_BAREMETAL
: Enable baremetal gathering, defaults to `true`

PROMETHEUS_SCW_ACCESS_KEY
: Access key for the Scaleway API

PROMETHEUS_SCW_SECRET_KEY
: Secret key for the Scaleway API

PROMETHEUS_SCW_ORG
: Organization for the Scaleway API

PROMETHEUS_SCW_ZONE
: Zone for the Scaleway API

PROMETHEUS_SCW_CONFIG
: Path to Scaleway configuration file

PROMETHEUS_SCW_INSTANCE_ZONES
: List of available zones for instance API, comma-separated list, defaults to `fr-par-1, nl-ams-1`

PROMETHEUS_SCW_BAREMETAL_ZONES
: List of available zones for baremetal API, comma-separated list, defaults to `fr-par-2`
