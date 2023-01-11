// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

module github.com/nginxinc/kubernetes-nginx-ingress

go 1.19

require github.com/sirupsen/logrus v1.9.0

require (
	github.com/nginxinc/nginx-plus-go-client v0.10.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/nginxinc/kubernetes-nginx-ingress/internal/config => ./internal/config
