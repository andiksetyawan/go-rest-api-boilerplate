FROM golang:1.18 as golang

RUN mkdir -p /

WORKDIR /

COPY . .

RUN make build