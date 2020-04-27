FROM golang:1.13.7-alpine3.11 as build

WORKDIR /app

ENV \
    TERM=xterm-color \
    TIME_ZONE="UTC" \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS="-mod=vendor" \
    GOLANGCI_VERSION="v1.23.3"

RUN \
    echo "## Prepare timezone" && \
    apk add --no-cache --update tzdata && \
    cp /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime && \
    echo "${TIME_ZONE}" > /etc/timezone && date && \
    echo "## Install golangci" && \
    wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VERSION} && \
    golangci-lint --version && \
    echo "## Prepare tdlib" && \
    apk add --no-cache gcc g++ zlib-dev openssl-dev libc6-compat && \
    apk add telegram-tdlib-static telegram-tdlib-dev --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted && \
    ln -s /usr/include /usr/local/include && \
    ln -s /usr/lib /usr/local/lib
ADD . .
RUN \
    go env && \
    go version && \
    echo "  ## Test" && \
    go test ./... && \
    echo "  ## Lint" && \
    golangci-lint run ./... && \
    echo "  ## Build" && \
    go build -o app . && \
    echo "  ## Done"

FROM golang:1.13.7-alpine3.11
RUN apk add telegram-tdlib-static telegram-tdlib-dev --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted && \
    ln -s /usr/include /usr/local/include && \
    ln -s /usr/lib /usr/local/lib && \
    mkdir -p /app/var/db && \
    mkdir -p /app/var/work && \
    chown -R nobody:nobody /app/var/
ENV \
    WORK_DIR=/app/var/work \
    PUBLIC_DIR=./public \
    DB=/app/var/db/groups.db
COPY --from=build /app/app /app/app
COPY --from=build /app/public /app/public
COPY --from=build /etc/localtime /etc/localtime
USER nobody:nobody
WORKDIR /app
CMD ./app