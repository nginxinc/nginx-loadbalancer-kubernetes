/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package certification

import (
	"context"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewCertificate(t *testing.T) {
	ctx := context.Background()

	cert, err := NewCertificates(ctx, nil)

	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if cert == nil {
		t.Fatalf(`cert should not be nil`)
	}
}

func TestCertificates_Initialize(t *testing.T) {
	certificates, err := buildCertificates()
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	err = certificates.Initialize()
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
}

func TestCertificates_RunWithoutInitialize(t *testing.T) {
	certificates, err := buildCertificates()
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	err = certificates.Run()
	if err == nil {
		t.Fatalf(`Expected error`)
	}

	if err.Error() != `initialize must be called before Run` {
		t.Fatalf(`Unexpected error: %v`, err)
	}
}

func TestCertificates_EmptyCertificates(t *testing.T) {
	certificates, err := buildCertificates()
	if err != nil {
		t.Fatalf(`error building Certificates: %v`, err)
	}

	err = certificates.Initialize()
	if err != nil {
		t.Fatalf(`error Initializing Certificates: %v`, err)
	}

	caBytes := certificates.GetCACertificate()
	if caBytes != nil {
		t.Fatalf(`Expected nil CA certificate`)
	}

	clientKey, clientCert := certificates.GetClientCertificate()
	if clientKey != nil {
		t.Fatalf(`Expected nil client key`)
	}
	if clientCert != nil {
		t.Fatalf(`Expected nil client certificate`)
	}
}

func buildCertificates() (*Certificates, error) {
	ctx := context.Background()
	k8sClient := fake.NewSimpleClientset()

	certificates, err := NewCertificates(ctx, k8sClient)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
