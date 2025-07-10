FROM --platform=$BUILDPLATFORM golang:1.24.5-alpine3.21@sha256:933e5a0829a1f97cc99917e3b799ebe450af30236f5a023a3583d26b5ef9166f AS builder

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

FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715

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
