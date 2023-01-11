// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/config"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/http"
	nginx "github.com/nginxinc/nginx-plus-go-client/client"
	"testing"
)

func TestNewNginxPlusSynchronizer(t *testing.T) {
	const NginxUrl = "http://demo.nginx.com/api"
	settings := config.Settings{NginxPlusHost: NginxUrl}
	httpClient, _ := http.NewHttpClient()
	nginxClient, _ := nginx.NewNginxClient(httpClient, settings.NginxPlusHost)

	synchronizer, err := NewNginxPlusSynchronizer(nginxClient)
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an NginxPlusSynchronizer instance")
	}
}
