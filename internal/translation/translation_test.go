// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package translation

//
//import (
//	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
//	v1 "k8s.io/api/networking/v1"
//	"testing"
//)
//
//func TestTranslateNoAddresses(t *testing.T) {
//	const expectedUpstreams = 0
//
//	ingress := &v1.Ingress{}
//	previousIngress := &v1.Ingress{}
//
//	event := core.NewEvent(core.Created, ingress, previousIngress)
//	updatedEvent, err := Translate(&event)
//	if err != nil {
//		t.Fatalf("Translate() error = %v", err)
//	}
//
//	actualUpstreams := len(updatedEvent.NginxUpstreams)
//	if actualUpstreams != expectedUpstreams {
//		t.Fatalf("expected %v upstreams, got %v", expectedUpstreams, actualUpstreams)
//	}
//}
//
//func TestTranslateOneIp(t *testing.T) {
//	const expectedUpstreams = 1
//
//	lbIngress := v1.IngressLoadBalancerIngress{
//		IP: "127.0.0.1",
//	}
//
//	ingress := &v1.Ingress{
//		Status: v1.IngressStatus{
//			LoadBalancer: v1.IngressLoadBalancerStatus{
//				Ingress: []v1.IngressLoadBalancerIngress{
//					lbIngress,
//				},
//			},
//		},
//	}
//
//	previousIngress := &v1.Ingress{}
//
//	event := core.NewEvent(core.Created, ingress, previousIngress)
//	updatedEvent, err := Translate(&event)
//	if err != nil {
//		t.Fatalf("Translate() error = %v", err)
//	}
//
//	actualUpstreams := len(updatedEvent.NginxUpstreams)
//	if actualUpstreams != expectedUpstreams {
//		t.Fatalf("expected %v upstreams, got %v", expectedUpstreams, actualUpstreams)
//	}
//}
//
//func TestTranslateOneHost(t *testing.T) {
//	const expectedUpstreams = 1
//
//	lbIngress := v1.IngressLoadBalancerIngress{
//		Hostname: "www.example.com",
//	}
//
//	ingress := &v1.Ingress{
//		Status: v1.IngressStatus{
//			LoadBalancer: v1.IngressLoadBalancerStatus{
//				Ingress: []v1.IngressLoadBalancerIngress{
//					lbIngress,
//				},
//			},
//		},
//	}
//
//	previousIngress := &v1.Ingress{}
//
//	event := core.NewEvent(core.Created, ingress, previousIngress)
//	updatedEvent, err := Translate(&event)
//	if err != nil {
//		t.Fatalf("Translate() error = %v", err)
//	}
//
//	actualUpstreams := len(updatedEvent.NginxUpstreams)
//	if actualUpstreams != expectedUpstreams {
//		t.Fatalf("expected %v upstreams, got %v", expectedUpstreams, actualUpstreams)
//	}
//}
//
//func TestTranslateOneHostAndOneIP(t *testing.T) {
//	const expectedUpstreams = 2
//
//	nodeIps := []string{"192.168.1.1"}
//
//	lbHostnameIngress := v1.IngressLoadBalancerIngress{
//		Hostname: "www.example.com",
//	}
//
//	lbIPIngress := v1.IngressLoadBalancerIngress{
//		IP: "127.0.0.1",
//	}
//
//	ingress := &v1.Ingress{
//		Status: v1.IngressStatus{
//			LoadBalancer: v1.IngressLoadBalancerStatus{
//				Ingress: []v1.IngressLoadBalancerIngress{
//					lbHostnameIngress,
//					lbIPIngress,
//				},
//			},
//		},
//	}
//
//	previousIngress := &v1.Ingress{}
//
//	event := core.NewEvent(core.Created, ingress, previousIngress, nodeIps)
//	updatedEvent, err := Translate(&event)
//	if err != nil {
//		t.Fatalf("Translate() error = %v", err)
//	}
//
//	actualUpstreams := len(updatedEvent.NginxUpstreams)
//	if actualUpstreams != expectedUpstreams {
//		t.Fatalf("expected %v upstreams, got %v", expectedUpstreams, actualUpstreams)
//	}
//}
//
//func TestTranslateMulitpleRoutes(t *testing.T) {
//	const expectedUpstreams = 6
//
//	lbHostnameIngress := v1.IngressLoadBalancerIngress{
//		Hostname: "www.example.com",
//	}
//
//	lbHostnameIngress1 := v1.IngressLoadBalancerIngress{
//		Hostname: "www.example.net",
//	}
//
//	lbHostnameIngress2 := v1.IngressLoadBalancerIngress{
//		Hostname: "www.example.org",
//	}
//
//	lbHostnameIngress3 := v1.IngressLoadBalancerIngress{
//		Hostname: "www.acme.com",
//	}
//
//	lbIPIngress := v1.IngressLoadBalancerIngress{
//		IP: "127.0.0.1",
//	}
//
//	lbIPIngress1 := v1.IngressLoadBalancerIngress{
//		IP: "192.168.0.1",
//	}
//
//	ingress := &v1.Ingress{
//		Status: v1.IngressStatus{
//			LoadBalancer: v1.IngressLoadBalancerStatus{
//				Ingress: []v1.IngressLoadBalancerIngress{
//					lbHostnameIngress,
//					lbHostnameIngress1,
//					lbHostnameIngress2,
//					lbHostnameIngress3,
//					lbIPIngress,
//					lbIPIngress1,
//				},
//			},
//		},
//	}
//
//	previousIngress := &v1.Ingress{}
//
//	event := core.NewEvent(core.Created, ingress, previousIngress)
//	updatedEvent, err := Translate(&event)
//	if err != nil {
//		t.Fatalf("Translate() error = %v", err)
//	}
//
//	actualUpstreams := len(updatedEvent.NginxUpstreams)
//	if actualUpstreams != expectedUpstreams {
//		t.Fatalf("expected %v upstreams, got %v", expectedUpstreams, actualUpstreams)
//	}
//}
