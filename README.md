# Prometheus Scaleway SD

[![Build Status](http://cloud.drone.io/api/badges/promhippie/prometheus-scw-sd/status.svg)](http://cloud.drone.io/promhippie/prometheus-scw-sd)
[![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4671e4dac861415db19d41c7959a530a)](https://www.codacy.com/app/promhippie/prometheus-scw-sd?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=promhippie/prometheus-scw-sd&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/promhippie/prometheus-scw-sd?status.svg)](http://godoc.org/github.com/promhippie/prometheus-scw-sd)
[![Go Report](http://goreportcard.com/badge/github.com/promhippie/prometheus-scw-sd)](http://goreportcard.com/report/github.com/promhippie/prometheus-scw-sd)
[![](https://images.microbadger.com/badges/image/promhippie/prometheus-scw-sd.svg)](http://microbadger.com/images/promhippie/prometheus-scw-sd "Get your own image badge on microbadger.com")

This project provides a server to automatically discover nodes within your Scaleway account in a Prometheus SD compatible format.

## Install

You can download prebuilt binaries from our [GitHub releases](https://github.com/promhippie/prometheus-scw-sd/releases), or you can use our Docker images published on [Docker Hub](https://hub.docker.com/r/promhippie/prometheus-scw-sd/tags/). If you need further guidance how to install this take a look at our [documentation](https://promhippie.github.io/prometheus-scw-sd/#getting-started).

## Integration

### Available labels

The following list of meta labels can be used to relabel your scrape results entirely. Hopefully the names are self-explaining, that's why I have skipped a description for each label.

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

### Prometheus config

Here you get a snippet for the Prometheus `scrape_config` that configures Prometheus to scrape `node_exporter` assuming that it is deployed on all your servers.

```
- job_name: node
  file_sd_configs:
  - files: [ "/etc/prometheus/scw.json" ]
  relabel_configs:
  - source_labels: [__meta_scaleway_public_ipv4]
    replacement: "${1}:9100"
    target_label: __address__
  - source_labels: [__meta_scaleway_zone]
    target_label: datacenter
  - source_labels: [__meta_scaleway_name]
    target_label: instance
```

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```bash
git clone https://github.com/promhippie/prometheus-scw-sd.git
cd prometheus-scw-sd

make generate build

./bin/prometheus-scw-sd -h
```

## Security

If you find a security issue please contact thomas@webhippie.de first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

* [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```
