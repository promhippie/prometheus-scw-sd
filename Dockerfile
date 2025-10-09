FROM --platform=$BUILDPLATFORM golang:1.25.2-alpine3.21@sha256:01346535ae797d5bc7301aa6518051e9a66adf813fc99e09872a06417759f913 AS builder

RUN apk add --no-cache -U git curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/prometheus-scw-sd
COPY . /go/src/prometheus-scw-sd/

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 prometheus-scw-sd && \
    adduser -D -u 1337 -h /var/lib/prometheus-scw-sd -G prometheus-scw-sd prometheus-scw-sd

EXPOSE 9000
VOLUME ["/var/lib/prometheus-scw-sd"]
ENTRYPOINT ["/usr/bin/prometheus-scw-sd"]
CMD ["server"]
HEALTHCHECK CMD ["/usr/bin/prometheus-scw-sd", "health"]

ENV PROMETHEUS_SCW_OUTPUT_ENGINE="http"
ENV PROMETHEUS_SCW_OUTPUT_FILE="/var/lib/prometheus-scw-sd/output.json"

COPY --from=builder /go/src/prometheus-scw-sd/bin/prometheus-scw-sd /usr/bin/prometheus-scw-sd
WORKDIR /var/lib/prometheus-scw-sd
USER prometheus-scw-sd
