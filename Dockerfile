FROM golang:1.19-alpine as build-env

ARG GOLANGCI_LINT_VERSION='v1.48.0'
ARG GOLANGCI_LINT_INSTALLER='https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh'

ARG SHELLCHECK_VERSION='v0.8.0'
ARG SHELLCHECK_TARBALL="https://github.com/koalaman/shellcheck/releases/download/${SHELLCHECK_VERSION?}/shellcheck-${SHELLCHECK_VERSION?}.linux.x86_64.tar.xz"

ARG UID=1000
ARG GID=100

RUN apk add --no-cache bash gcc make musl-dev
RUN addgroup -g $GID devgroup || true && \
    adduser -D -u $UID -G "$(getent group 20 | cut -d ':' -f1)" dev

ENV GOCACHE="/home/dev/go/cache"
ENV GOMODCACHE="/home/dev/go/pkg/mod"

RUN mkdir -p "$GOCACHE" && mkdir -p "$GOMODCACHE"
RUN wget -O- -nv "$GOLANGCI_LINT_INSTALLER" | sh -s "$GOLANGCI_LINT_VERSION" -b /usr/local/bin &&  \
    \
    wget -qO- "$SHELLCHECK_TARBALL" | tar -xJv && \
    cp "shellcheck-${SHELLCHECK_VERSION}/shellcheck" /usr/local/bin/ &&  \
    \
    rm -rf "shellcheck-${SHELLCHECK_VERSION}"

RUN chown -R $UID:$GID /home/dev/

VOLUME "$GOCACHE"
VOLUME "$GOMODCACHE"

USER dev

WORKDIR /usr/src/app
COPY . .

RUN make DOCKER_ENABLED=0 all

CMD exec /bin/sh -c 'trap : TERM INT; sleep infinity & wait'


FROM alpine

COPY --from=build-env /usr/src/app/bin/benchmark /usr/local/bin/
WORKDIR /data

ENTRYPOINT ["benchmark"]
