FROM golang:1.18.5-alpine AS base

ENV GOPATH=/opt/service/.go

WORKDIR /opt/service/
COPY . .

RUN go mod download
RUN go build -o /bin/api ./cmd/api

ENTRYPOINT ["/bin/api"]