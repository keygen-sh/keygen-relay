FROM golang:1.23-alpine AS build
WORKDIR /go/src/github.com/keygen-sh/keygen-relay
RUN apk add --no-cache gcc musl-dev sqlite-dev make
COPY . .
RUN make build

FROM alpine:latest
RUN apk add --no-cache bash openssh-client sqlite-libs
COPY --from=build /go/src/github.com/keygen-sh/keygen-relay/bin/* /usr/bin/
VOLUME /app/data
ENV RELAY_DATABASE=/app/data/relay.sqlite
WORKDIR /app

ENTRYPOINT ["relay"]
