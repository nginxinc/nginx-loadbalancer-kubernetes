# Copyright 2023 f5 Inc. All rights reserved.
# Use of this source code is governed by the Apache
# license that can be found in the LICENSE file.

FROM golang:1.19.5-alpine3.16 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o nginx-k8s-edge-controller ./cmd/nginx-k8s-edge-controller/main.go

FROM alpine:3.16

WORKDIR /opt/nginx-k8s-edge-controller

COPY --from=builder /app/nginx-k8s-edge-controller .

ENTRYPOINT ["/opt/nginx-k8s-edge-controller/nginx-k8s-edge-controller"]
