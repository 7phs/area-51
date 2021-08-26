FROM golang:1.17-buster

ENV SRC=/area-51

ADD . /src
WORKDIR /src

RUN make build

FROM debian:stretch

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get clean

WORKDIR /root/
COPY --from=0 ${SRC}/bin ./app

CMD ["./app/server"]