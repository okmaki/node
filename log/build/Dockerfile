FROM golang:latest AS build

RUN mkdir /build
WORKDIR /build
RUN git clone https://github.com/okmaki/node.git
WORKDIR /build/node/log
RUN go build ./cmd/service

FROM debian:latest AS run

RUN mkdir -p /app
COPY --from=build /build/node/log/service /app
ENTRYPOINT [ "/app/log" ]
