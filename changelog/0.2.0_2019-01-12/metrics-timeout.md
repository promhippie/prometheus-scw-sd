Change: Timeout for metrics handler

We added an additional middleware to properly timeout requests to the metrics
endpoint for long running request.

https://github.com/promhippie/prometheus-scw-sd/pull/3
