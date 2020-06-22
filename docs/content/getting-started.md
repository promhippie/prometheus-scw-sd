---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 10
---

## Installation

We won't cover further details how to properly setup [Prometheus](https://prometheus.io) itself, we will only cover some basic setup based on [docker-compose](https://docs.docker.com/compose/). But if you want to run this service discovery without [docker-compose](https://docs.docker.com/compose/) you should be able to adopt that to your needs.

First of all we need to prepare a configuration for [Prometheus](https://prometheus.io) that includes the service discovery which simply maps to a node exporter.

{{< highlight yaml >}}
global:
  scrape_interval: 1m
  scrape_timeout: 10s
  evaluation_interval: 1m

scrape_configs:
- job_name: node
  file_sd_configs:
  - files: [ "/etc/sd/scw.json" ]
  relabel_configs:
  - source_labels: [__meta_scaleway_public_ipv4]
    replacement: "${1}:9100"
    target_label: __address__
  - source_labels: [__meta_scaleway_zone]
    target_label: datacenter
  - source_labels: [__meta_scaleway_name]
    target_label: instance
- job_name: scw-sd
  static_configs:
  - targets:
    - scw-sd:9000
{{< / highlight >}}

After preparing the configuration we need to create the `docker-compose.yml` within the same folder, this `docker-compose.yml` starts a simple [Prometheus](https://prometheus.io) instance together with the service discovery. Don't forget to update the envrionment variables with the required credentials. If you are using a different volume for the service discovery you have to make sure that the container user is allowed to write to this volume.

{{< highlight yaml >}}
version: '2'

volumes:
  prometheus:

services:
  prometheus:
    image: prom/prometheus:v2.6.0
    restart: always
    ports:
      - 9090:9090
    volumes:
      - prometheus:/prometheus
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./service-discovery:/etc/sd

  scw-sd:
    image: promhippie/prometheus-scw-sd:latest
    restart: always
    environment:
      - PROMETHEUS_SCW_LOG_PRETTY=true
      - PROMETHEUS_SCW_OUTPUT_FILE=/etc/sd/scw.json
      - PROMETHEUS_SCW_ACCESS_KEY=SCWTAGNAXNNTV4NBNPVN
      - PROMETHEUS_SCW_SECRET_KEY=a6566958-0311-4c72-bac0-2cc86a852a19
      - PROMETHEUS_SCW_ORG=a84ef57f-c727-40c4-a236-1a81e8041ce5
      - PROMETHEUS_SCW_REGION=par1
    volumes:
      - ./service-discovery:/etc/sd
{{< / highlight >}}

Since our `latest` Docker tag always refers to the `master` branch of the Git repository you should always use some fixed version. You can see all available tags at our [DockerHub repository](https://hub.docker.com/r/promhippie/prometheus-scw-sd/tags/), there you will see that we also provide a manifest, you can easily start the exporter on various architectures without any change to the image name. You should apply a change like this to the `docker-compose.yml`:

{{< highlight diff >}}
  scw-sd:
-   image: promhippie/prometheus-scw-sd:latest
+   image: promhippie/prometheus-scw-sd:0.2.0
    restart: always
    environment:
      - PROMETHEUS_SCW_LOG_PRETTY=true
      - PROMETHEUS_SCW_OUTPUT_FILE=/etc/sd/scw.json
      - PROMETHEUS_SCW_ACCESS_KEY=SCWTAGNAXNNTV4NBNPVN
      - PROMETHEUS_SCW_SECRET_KEY=a6566958-0311-4c72-bac0-2cc86a852a19
      - PROMETHEUS_SCW_ORG=a84ef57f-c727-40c4-a236-1a81e8041ce5
      - PROMETHEUS_SCW_REGION=par1
    volumes:
      - ./service-discovery:/etc/sd
{{< / highlight >}}

Depending on how you have launched and configured [Prometheus](https://prometheus.io) it's possible that it's running as user `nobody`, in that case you should run the service discovery as this user as well, otherwise [Prometheus](https://prometheus.io) won't be able to read the generated JSON file:

{{< highlight diff >}}
  scw-sd:
    image: promhippie/prometheus-scw-sd:latest
    restart: always
+   user: '65534'
    environment:
      - PROMETHEUS_SCW_LOG_PRETTY=true
      - PROMETHEUS_SCW_OUTPUT_FILE=/etc/sd/scw.json
      - PROMETHEUS_SCW_ACCESS_KEY=SCWTAGNAXNNTV4NBNPVN
      - PROMETHEUS_SCW_SECRET_KEY=a6566958-0311-4c72-bac0-2cc86a852a19
      - PROMETHEUS_SCW_ORG=a84ef57f-c727-40c4-a236-1a81e8041ce5
      - PROMETHEUS_SCW_REGION=par1
    volumes:
      - ./service-discovery:/etc/sd
{{< / highlight >}}

Finally the service discovery should be configured fine, let's start this stack with [docker-compose](https://docs.docker.com/compose/), you just need to execute `docker-compose up` within the directory where you have stored `prometheus.yml` and `docker-compose.yml`.

{{< highlight txt >}}
# docker-compose up
Creating network "scw-sd_default" with the default driver
Creating volume "scw-sd_prometheus" with default driver
Creating scw-sd_scw-sd_1     ... done
Creating scw-sd_prometheus_1 ... done
Attaching to scw-sd_scw-sd_1, scw-sd_prometheus_1
prometheus_1  | level=info ts=2018-10-07T16:18:45.1092588Z caller=main.go:238 msg="Starting Prometheus" version="(version=2.4.3, branch=HEAD, revision=167a4b4e73a8eca8df648d2d2043e21bdb9a7449)"
prometheus_1  | level=info ts=2018-10-07T16:18:45.1093473Z caller=main.go:239 build_context="(go=go1.11.1, user=root@1e42b46043e9, date=20181004-08:42:02)"
prometheus_1  | level=info ts=2018-10-07T16:18:45.1093727Z caller=main.go:240 host_details="(Linux 4.9.93-linuxkit-aufs #1 SMP Wed Jun 6 16:55:56 UTC 2018 x86_64 a13beb141445 (none))"
prometheus_1  | level=info ts=2018-10-07T16:18:45.1093944Z caller=main.go:241 fd_limits="(soft=1048576, hard=1048576)"
prometheus_1  | level=info ts=2018-10-07T16:18:45.1094111Z caller=main.go:242 vm_limits="(soft=unlimited, hard=unlimited)"
prometheus_1  | level=info ts=2018-10-07T16:18:45.1103691Z caller=main.go:554 msg="Starting TSDB ..."
prometheus_1  | level=info ts=2018-10-07T16:18:45.1104236Z caller=web.go:397 component=web msg="Start listening for connections" address=0.0.0.0:9090
prometheus_1  | level=info ts=2018-10-07T16:18:45.1162682Z caller=main.go:564 msg="TSDB started"
prometheus_1  | level=info ts=2018-10-07T16:18:45.1171396Z caller=main.go:624 msg="Loading configuration file" filename=/etc/prometheus/prometheus.yml
prometheus_1  | level=info ts=2018-10-07T16:18:45.11963Z caller=main.go:650 msg="Completed loading of configuration file" filename=/etc/prometheus/prometheus.yml
prometheus_1  | level=info ts=2018-10-07T16:18:45.1196724Z caller=main.go:523 msg="Server is ready to receive web requests."
scw-sd_1      | level=info ts=2018-10-07T16:18:45.1296711Z msg="Launching Prometheus Scaleway SD" version=0.0.0-master revision=a94cb2f date=20180924 go=go1.11
scw-sd_1      | level=info ts=2018-10-07T16:18:45.8562717Z msg="Starting metrics server" addr=0.0.0.0:9000
{{< / highlight >}}

That's all, the service discovery should be up and running. You can access [Prometheus](https://prometheus.io) at [http://localhost:9090](http://localhost:9090).

{{< figure src="service-discovery.png" title="Prometheus service discovery for Hetzner" >}}

## Kubernetes

Currently we have not prepared a deployment for Kubernetes, but this is something we will provide for sure. Most interesting will be the integration into the [Prometheus Operator](https://coreos.com/operators/prometheus/docs/latest/), so stay tuned.

## Configuration

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below, in case you want to configure multiple accounts with a single service you are forced to use the configuration file as the environment variables are limited to a single account. As the service is pretty lightweight you can even start an instance per account and configure it entirely by the variables, it's up to you.

PROMETHEUS_SCW_CONFIG
: Path to Scaleway configuration file, optionally, required for multi credentials

PROMETHEUS_SCW_ACCESS_KEY
: Access key for the Scaleway API, required for authentication

PROMETHEUS_SCW_SECRET_KEY
: Secret key for the Scaleway API, required for authentication

PROMETHEUS_SCW_ORG
: Organization for the Scaleway API, optionally

PROMETHEUS_SCW_ZONE
: Zone for the Scaleway API, optionally

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

PROMETHEUS_SCW_CHECK_INSTANCE
: Enable fetching servers from instance API, defaults to `true`

PROMETHEUS_SCW_CHECK_BAREMETAL
: Enable fetching servers from baremetal API, defaults to `true`

### Configuration file

Especially if you want to configure multiple accounts within a single service discovery you got to use the configuration file. So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/promhippie/prometheus-scw-sd/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options, they also include the default values.

## Labels

* `__address__`
* `__meta_scaleway_allowed_actions`
* `__meta_scaleway_arch`
* `__meta_scaleway_boot_type`
* `__meta_scaleway_bootscript_id`
* `__meta_scaleway_bootscript_initrd`
* `__meta_scaleway_bootscript_kernel`
* `__meta_scaleway_bootscript_title`
* `__meta_scaleway_cluster`
* `__meta_scaleway_commercial_type`
* `__meta_scaleway_description`
* `__meta_scaleway_domain`
* `__meta_scaleway_dynamic_ip_required`
* `__meta_scaleway_enable_ipv6`
* `__meta_scaleway_hostname`
* `__meta_scaleway_hypervisor`
* `__meta_scaleway_id`
* `__meta_scaleway_image_id`
* `__meta_scaleway_image_name`
* `__meta_scaleway_install_hostname`
* `__meta_scaleway_install_os`
* `__meta_scaleway_install_status`
* `__meta_scaleway_ips`
* `__meta_scaleway_ipv6`
* `__meta_scaleway_kind`
* `__meta_scaleway_name`
* `__meta_scaleway_node`
* `__meta_scaleway_offer`
* `__meta_scaleway_org`
* `__meta_scaleway_placement_group_id`
* `__meta_scaleway_placement_group_name`
* `__meta_scaleway_platform`
* `__meta_scaleway_private_host`
* `__meta_scaleway_private_ipv4`
* `__meta_scaleway_project`
* `__meta_scaleway_protected`
* `__meta_scaleway_public_host`
* `__meta_scaleway_public_ipv4`
* `__meta_scaleway_security_group_id`
* `__meta_scaleway_security_group_name`
* `__meta_scaleway_state_detail`
* `__meta_scaleway_state`
* `__meta_scaleway_status`
* `__meta_scaleway_tags`
* `__meta_scaleway_zone`

## Metrics

prometheus_scw_sd_request_duration_seconds
: Histogram of latencies for requests to the Scaleway API

prometheus_scw_sd_request_failures_total
: Total number of failed requests to the Scaleway API
