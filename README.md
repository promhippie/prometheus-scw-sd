# prometheus-scw-sd

[![Current Tag](https://img.shields.io/github/v/tag/promhippie/prometheus-scw-sd?sort=semver)](https://github.com/promhippie/prometheus-scw-sd) [![Build Status](https://github.com/promhippie/prometheus-scw-sd/actions/workflows/general.yml/badge.svg)](https://github.com/promhippie/prometheus-scw-sd/actions) [![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#promhippie:matrix.org) [![Docker Size](https://img.shields.io/docker/image-size/promhippie/prometheus-scw-sd/latest)](https://hub.docker.com/r/promhippie/prometheus-scw-sd) [![Docker Pulls](https://img.shields.io/docker/pulls/promhippie/prometheus-scw-sd)](https://hub.docker.com/r/promhippie/prometheus-scw-sd) [![Go Reference](https://pkg.go.dev/badge/github.com/promhippie/prometheus-scw-sd.svg)](https://pkg.go.dev/github.com/promhippie/prometheus-scw-sd) [![Go Report Card](https://goreportcard.com/badge/github.com/promhippie/prometheus-scw-sd)](https://goreportcard.com/report/github.com/promhippie/prometheus-scw-sd) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/4671e4dac861415db19d41c7959a530a)](https://www.codacy.com/gh/promhippie/prometheus-scw-sd/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=promhippie/prometheus-scw-sd&amp;utm_campaign=Badge_Grade)

This project provides a server to automatically discover nodes within your
Scaleway account in a Prometheus SD compatible format.

## Install

You can download prebuilt binaries from our [GitHub releases][releases]. Beside
that we are publishing Docker images to [Docker Hub][dockerhub] and
[Quay][quay]. If you need further guidance how to install this take a look at
our [documentation][docs].

## Development

Make sure you have a working Go environment, for further reference or a guide
take a look at the [install instructions][golang]. This project requires
Go >= v1.17, at least that's the version we are using.

```console
git clone https://github.com/promhippie/prometheus-scw-sd.git
cd prometheus-scw-sd

make generate build

./bin/prometheus-scw-sd -h
```

## Security

If you find a security issue please contact
[thomas@webhippie.de](mailto:thomas@webhippie.de) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

-   [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2018 Thomas Boerger <thomas@webhippie.de>
```

[releases]: https://github.com/promhippie/prometheus-scw-sd/releases
[dockerhub]: https://hub.docker.com/r/promhippie/prometheus-scw-sd/tags/
[quay]: https://quay.io/repository/promhippie/prometheus-scw-sd?tab=tags
[docs]: https://promhippie.github.io/prometheus-scw-sd/#getting-started
[golang]: http://golang.org/doc/install.html
