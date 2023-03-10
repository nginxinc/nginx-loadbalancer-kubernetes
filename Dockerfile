# Copyright 2023 f5 Inc. All rights reserved.
# Use of this source code is governed by the Apache
# license that can be found in the LICENSE file.

FROM golang:1.19.5-alpine3.16 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o nginx-k8s-loadbalancer ./cmd/nginx-k8s-loadbalancer/main.go

FROM alpine:3.16

WORKDIR /opt/nginx-k8s-loadbalancer

COPY --from=builder /app/nginx-k8s-loadbalancer .

ENTRYPOINT ["/opt/nginx-k8s-loadbalancer/nginx-k8s-loadbalancer"]
