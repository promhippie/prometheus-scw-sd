---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 10
---

## Installation

We won't cover further details how to properly setup [Prometheus](https://prometheus.io) itself, we will only cover some basic setup based on [docker-compose](https://docs.docker.com/compose/). But if you want to run this service discovery without [docker-compose](https://docs.docker.com/compose/) you should be able to adopt that to your needs.

First of all we need to prepare a configuration for [Prometheus](https://prometheus.io) that includes the service discovery which simply maps to a node exporter.

{{< gist tboerger 4f445d609f0894a5fe40c20e7a42d565 "prometheus.yml" >}}

After preparing the configuration we need to create the `docker-compose.yml` within the same folder, this `docker-compose.yml` starts a simple [Prometheus](https://prometheus.io) instance together with the service discovery. Don't forget to update the envrionment variables with the required credentials. If you are using a different volume for the service discovery you have to make sure that the container user is allowed to write to this volume.

{{< gist tboerger 4f445d609f0894a5fe40c20e7a42d565 "docker-compose.yml" >}}

Since our `latest` Docker tag always refers to the `master` branch of the Git repository you should always use some fixed version. You can see all available tags at our [DockerHub repository](https://hub.docker.com/r/promhippie/prometheus-scw-sd/tags/), there you will see that we also provide a manifest, you can easily start the exporter on various architectures without any change to the image name. You should apply a change like this to the `docker-compose.yml`:

{{< gist tboerger 4f445d609f0894a5fe40c20e7a42d565 "tag.diff" >}}

Depending on how you have launched and configured [Prometheus](https://prometheus.io) it's possible that it's running as user `nobody`, in that case you should run the service discovery as this user as well, otherwise [Prometheus](https://prometheus.io) won't be able to read the generated JSON file:

{{< gist tboerger 4f445d609f0894a5fe40c20e7a42d565 "userid.diff" >}}

Finally the service discovery should be configured fine, let's start this stack with [docker-compose](https://docs.docker.com/compose/), you just need to execute `docker-compose up` within the directory where you have stored `prometheus.yml` and `docker-compose.yml`.

{{< gist tboerger 4f445d609f0894a5fe40c20e7a42d565 "output.log" >}}

That's all, the service discovery should be up and running. You can access [Prometheus](https://prometheus.io) at [http://localhost:9090](http://localhost:9090).

{{< figure src="service-discovery.png" title="Prometheus service discovery for Hetzner" >}}

## Kubernetes

Currently we have not prepared a deployment for Kubernetes, but this is something we will provide for sure. Most interesting will be the integration into the [Prometheus Operator](https://coreos.com/operators/prometheus/docs/latest/), so stay tuned.

## Configuration

PROMETHEUS_SCW_TOKEN
: Access token for the Scaleway API, required for authentication

PROMETHEUS_SCW_ORG
: Organization for the Scaleway API, required for authentication

PROMETHEUS_SCW_REGION
: Region for the Scaleway API, required for authentication

PROMETHEUS_SCW_LOG_LEVEL
: Only log messages with given severity, defaults to `info`

PROMETHEUS_SCW_LOG_PRETTY
: Enable pretty messages for logging, defaults to `true`

PROMETHEUS_SCW_WEB_ADDRESS
: Address to bind the metrics server, defaults to `0.0.0.0:9000`

PROMETHEUS_SCW_WEB_PATH
: Path to bind the metrics server, defaults to `/metrics`

PROMETHEUS_SCW_OUTPUT_FILE
: Path to write the file_sd config, defaults to `/etc/prometheus/scw.json`

PROMETHEUS_SCW_OUTPUT_REFRESH
: Discovery refresh interval in seconds, defaults to `30`

## Labels

* `__meta_scaleway_name`
* `__meta_scaleway_id`
* `__meta_scaleway_arch`
* `__meta_scaleway_image_id`
* `__meta_scaleway_image_name`
* `__meta_scaleway_public_ipv4`
* `__meta_scaleway_public_host`
* `__meta_scaleway_state`
* `__meta_scaleway_private_ipv4`
* `__meta_scaleway_private_host`
* `__meta_scaleway_hostname`
* `__meta_scaleway_org`
* `__meta_scaleway_commercial_type`
* `__meta_scaleway_platform`
* `__meta_scaleway_hypervisor`
* `__meta_scaleway_node`
* `__meta_scaleway_blade`
* `__meta_scaleway_chassis`
* `__meta_scaleway_cluster`
* `__meta_scaleway_zone`
* `__meta_scaleway_tags`

## Metrics

prometheus_scw_sd_request_duration_seconds
: Histogram of latencies for requests to the Scaleway API

prometheus_scw_sd_request_failures_total
: Total number of failed requests to the Scaleway API
