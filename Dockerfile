# Copyright 2023 f5 Inc. All rights reserved.
# Use of this source code is governed by the Apache
# license that can be found in the LICENSE file.

FROM golang:1.19.5-alpine3.16 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o nginx-loadbalancer-kubernetes ./cmd/nginx-loadbalancer-kubernetes/main.go

FROM alpine:3.16

WORKDIR /opt/nginx-loadbalancer-kubernetes

RUN adduser -u 11115 -D -H  nlk

USER nlk

COPY --from=builder /app/nginx-loadbalancer-kubernetes .

ENTRYPOINT ["/opt/nginx-loadbalancer-kubernetes/nginx-loadbalancer-kubernetes"]
