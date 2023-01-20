// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package nginxplus

import (
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	v1 "k8s.io/api/networking/v1"
	"testing"
)

func TestTranslateUnknown(t *testing.T) {
	ingress := &v1.Ingress{}
	var previousIngress *v1.Ingress

	event := core.NewEvent(-1, ingress, previousIngress)
	_, err := Translate(&event)
	if err.Error() != "unknown event type" {
		t.Fatalf("Expected an error %v", err)
	}
}

func TestTranslateCreated(t *testing.T) {
	ingress := &v1.Ingress{}
	previousIngress := &v1.Ingress{}

	event := core.NewEvent(core.Created, ingress, previousIngress)
	_, err := Translate(&event)
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
}

func TestTranslateUpdated(t *testing.T) {
	ingress := &v1.Ingress{}
	previousIngress := &v1.Ingress{}

	event := core.NewEvent(core.Updated, ingress, previousIngress)
	_, err := Translate(&event)
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
}

func TestTranslateDeleted(t *testing.T) {
	ingress := &v1.Ingress{}
	var previousIngress *v1.Ingress

	event := core.NewEvent(core.Deleted, ingress, previousIngress)
	_, err := Translate(&event)
	if err != nil {
		t.Fatalf("Translate() error = %v", err)
	}
}
