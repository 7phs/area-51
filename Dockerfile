FROM golang:1.17-buster as builder

ENV SRC=/area-51

ADD . /src
WORKDIR /src

RUN go build -o /bin/ /src/cmd/...

FROM debian:stretch

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get clean

COPY --from=builder /bin /app

ENTRYPOINT ["/app/server"]