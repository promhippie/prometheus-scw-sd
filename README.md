# Prometheus Scaleway SD

[![Current Tag](https://img.shields.io/github/v/tag/promhippie/prometheus-scw-sd?sort=semver)](https://github.com/promhippie/prometheus-scw-sd) [![General Build](https://github.com/promhippie/prometheus-scw-sd/actions/workflows/general.yml/badge.svg)](https://github.com/promhippie/prometheus-scw-sd/actions/workflows/general.yaml) [![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/87cbb93f28be43a2a871018f106bc286)](https://www.codacy.com/gh/promhippie/prometheus-scw-sd/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=promhippie/prometheus-scw-sd&amp;utm_campaign=Badge_Grade) [![Go Doc](https://godoc.org/github.com/promhippie/prometheus-scw-sd?status.svg)](http://godoc.org/github.com/promhippie/prometheus-scw-sd) [![Go Report](http://goreportcard.com/badge/github.com/promhippie/prometheus-scw-sd)](http://goreportcard.com/report/github.com/promhippie/prometheus-scw-sd)

This project provides a server to automatically discover nodes within your
Scaleway account in a Prometheus SD compatible format.

## Install

You can download prebuilt binaries from our [GitHub releases][releases]. Besides
that we also prepared repositories for DEB and RPM packages which can be  found
at [Baltorepo][baltorepo]. If you prefer to use containers you could use our
images published on [GHCR][ghcr], [Docker Hub][dockerhub] or [Quay][quayio]. If
you need further guidance how to install this take a look at our [docs][docs].

## Development

If you are not familiar with [Nix][nix] it is up to you to have a working
environment for Go (>= 1.24.0) as the setup won't we covered within this guide.
Please follow the official install instructions for [Go][golang]. Beside that
we are using [go-task][gotask] to define all commands to build this project.

```console
git clone https://github.com/promhippie/prometheus-scw-sd.git
cd prometheus-scw-sd

task generate build
./bin/prometheus-scw-sd -h
```

If you got [Nix][nix] and [Direnv][direnv] configured you can simply execute
the following commands to get al dependencies including [go-task][gotask] and
the required runtimes installed. You are also able to directly use the process
manager of [devenv][devenv]:

```console
cat << EOF > .envrc
use flake . --impure --extra-experimental-features nix-command
EOF

direnv allow
```

To start developing on this project you have to execute only a few commands:

```console
task watch
```

The development server should be running on
[http://localhost:9000](http://localhost:9000). Generally it supports
hot reloading which means the services are automatically restarted/reloaded on
code changes.

If you got [Nix][nix] configured you can simply execute the [devenv][devenv]
command to start:

```console
devenv up
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
[baltorepo]: https://webhippie.baltorepo.com/promhippie/
[ghcr]: https://github.com/promhippie/prometheus-scw-sd/pkgs/container/prometheus-scw-sd
[dockerhub]: https://hub.docker.com/r/promhippie/prometheus-scw-sd/tags/
[quayio]: https://quay.io/repository/promhippie/prometheus-scw-sd?tab=tags
[docs]: https://promhippie.github.io/prometheus-scw-sd/#getting-started
[nix]: https://nixos.org/
[golang]: http://golang.org/doc/install.html
[gotask]: https://taskfile.dev/installation/
[direnv]: https://direnv.net/
[devenv]: https://devenv.sh/
