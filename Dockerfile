FROM golang:1.13.7-alpine3.11 as build
ENV \
    MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/v4.10.0/migrate.linux-amd64.tar.gz" \
    PG_DATA_DIR="/var/lib/postgresql/data" \
    PG_DSN="postgres://postgres:@localhost/test?sslmode=disable" \
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
    echo "${TIME_ZONE}" > /etc/timezone && date
RUN \
    echo "## Install PostgreSQL" && \
    apk add --no-cache --update postgresql postgresql-client postgresql-contrib tar curl && \
    mkdir -p /run/postgresql/ && \
    chown -R postgres /run/postgresql
RUN \
    echo "## Install migrate" && \
    curl --location --silent ${MIGRATE_URL} | tar -xvz -C /tmp/ && \
    mv /tmp/migrate.linux-amd64 /bin/migrate && \
    migrate -version
RUN \
    echo "## Setup PostgreSQL" && \
    mkdir -p /var/lib/postgresql && \
    su postgres -c 'pg_ctl initdb --pgdata ${PG_DATA_DIR} -U postgres'
RUN \
    echo "## Prepare tdlib" && \
    apk add --no-cache gcc g++ zlib-dev openssl-dev libc6-compat && \
    apk add telegram-tdlib-static telegram-tdlib-dev --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted && \
    ln -s /usr/include /usr/local/include && \
    ln -s /usr/lib /usr/local/lib
RUN \
    echo "## Install golangci" && \
    wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VERSION} && \
    golangci-lint --version

WORKDIR /app

ADD ./db /db
RUN \
    echo "## Run PostgreSQL" && \
    su postgres -c 'pg_ctl start --pgdata ${PG_DATA_DIR}' && \
    psql -U postgres -h localhost -c "CREATE DATABASE test" && \
    echo "## Test migrations" && \
    echo "    ## Migrate up" && \
    migrate -path /db -database ${PG_DSN} -verbose up && \
    echo "    ## Migrate down" && \
    migrate -path /db -database ${PG_DSN} -verbose goto 1 && \
    echo "    ## Migrate up" && \
    migrate -path /db -database ${PG_DSN} -verbose up

ADD . .
RUN \
    echo "## Test" && \
    go test ./... && \
    echo "## Lint" && \
    golangci-lint run ./... && \
    echo "## Build" && \
    go build -o app . && \
    echo "## Done"

FROM alpine:3.11
ENV \
    TMP_DIR=/app/tmp \
    PUBLIC_DIR=/app/public
COPY --from=build /etc/localtime /etc/localtime
COPY --from=build /app/app /app/app
COPY --from=build /app/public /app/public
COPY --from=build /app/db  /migrations
COPY --from=build /bin/migrate  /bin/migrate
RUN set -ex \
    apk add --no-cache gcc g++ zlib-dev openssl-dev libc6-compat && \
    apk add telegram-tdlib-static telegram-tdlib-dev --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ --allow-untrusted && \
    rm -f /var/cache/apk/* && \
    ln -s /usr/include /usr/local/include && \
    rm -rf /usr/local/lib && \
    ln -s /usr/lib /usr/local/lib && \
    mkdir -p ${TMP_DIR} && \
    chown -R nobody:nobody ${TMP_DIR}
USER nobody:nobody
CMD /app/app