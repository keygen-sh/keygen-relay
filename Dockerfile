FROM golang:1.23-alpine AS base
FROM sqlc/sqlc:latest AS sqlc

FROM base AS build
WORKDIR /go/src/github.com/keygen-sh/keygen-relay
RUN apk add --no-cache gcc musl-dev sqlite-dev make
COPY --from=sqlc /workspace/sqlc /usr/bin/sqlc
COPY . .
RUN make build

FROM alpine:latest
RUN apk add --no-cache bash openssh-client sqlite-libs
COPY --from=build /go/src/github.com/keygen-sh/keygen-relay/bin/* /usr/bin/

VOLUME /data
WORKDIR /data

ENV RELAY_DATABASE=/data/relay.sqlite
ENV RELAY_ADDR=0.0.0.0
ENV RELAY_PORT=6349

ENTRYPOINT ["relay"]

EXPOSE 6349
CMD ["serve"]
