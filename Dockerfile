FROM golang:1.18 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN make build-linux

FROM debian:stable-slim

COPY --from=builder /app/bin/si-linux-amd64 /usr/local/bin/si

RUN useradd -r -d / -s /sbin/nologin -c 'Application User' -u 1001 -g 0 app && \
    chmod +x /usr/local/bin/si

USER 1001

ENTRYPOINT ["si"]
