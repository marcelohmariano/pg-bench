ARG WORKDIR="/usr/src/app"

FROM golang:1.19-alpine as builder

ARG WORKDIR

WORKDIR "$WORKDIR"
COPY . .

RUN apk add --no-cache bash make && make DOCKER_ENABLED=0 all


FROM alpine as release

ARG WORKDIR

COPY --from=builder "${WORKDIR}/bin/benchmark" /usr/local/bin/benchmark
WORKDIR /data

ENTRYPOINT ["benchmark"]
CMD ["--help"]


FROM golang:1.19 as build-env

ARG GOLANGCI_LINT_VERSION='v1.48.0'
ARG GOLANGCI_LINT_INSTALLER='https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh'

ARG UID=1000
ARG GID=100

RUN apt update && apt install -y shellcheck
RUN groupadd -g $GID devgroup || true && \
    useradd -m -u $UID -g "$(getent group $GID | cut -d ':' -f1)" dev

ENV GOCACHE="/home/dev/go/cache"
ENV GOMODCACHE="/home/dev/go/pkg/mod"

RUN mkdir -p "$GOCACHE" && mkdir -p "$GOMODCACHE"
RUN wget -O- -nv "$GOLANGCI_LINT_INSTALLER" | sh -s "$GOLANGCI_LINT_VERSION" -b /usr/local/bin

RUN chown -R $UID:$GID /home/dev/

VOLUME "$GOCACHE"
VOLUME "$GOMODCACHE"

WORKDIR "$WORKDIR"
USER dev

CMD exec /bin/sh -c 'trap : TERM INT; sleep infinity & wait'
