/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package configuration

import (
	"testing"
)

func Test_String(t *testing.T) {
	t.Parallel()
	mode := NoTLS.String()
	if mode != "no-tls" {
		t.Errorf("Expected TLSModeNoTLS to be 'no-tls', got '%s'", mode)
	}

	mode = CertificateAuthorityTLS.String()
	if mode != "ca-tls" {
		t.Errorf("Expected TLSModeCaTLS to be 'ca-tls', got '%s'", mode)
	}

	mode = CertificateAuthorityMutualTLS.String()
	if mode != "ca-mtls" {
		t.Errorf("Expected TLSModeCaMTLS to be 'ca-mtls', got '%s'", mode)
	}

	mode = SelfSignedTLS.String()
	if mode != "ss-tls" {
		t.Errorf("Expected TLSModeSsTLS to be 'ss-tls', got '%s'", mode)
	}

	mode = SelfSignedMutualTLS.String()
	if mode != "ss-mtls" {
		t.Errorf("Expected TLSModeSsMTLS to be 'ss-mtls', got '%s',", mode)
	}

	mode = TLSMode(5).String()
	if mode != "" {
		t.Errorf("Expected TLSMode(5) to be '', got '%s'", mode)
	}
}

func Test_TLSModeMap(t *testing.T) {
	t.Parallel()
	mode := TLSModeMap["no-tls"]
	if mode != NoTLS {
		t.Errorf("Expected TLSModeMap['no-tls'] to be TLSModeNoTLS, got '%d'", mode)
	}

	mode = TLSModeMap["ca-tls"]
	if mode != CertificateAuthorityTLS {
		t.Errorf("Expected TLSModeMap['ca-tls'] to be TLSModeCaTLS, got '%d'", mode)
	}

	mode = TLSModeMap["ca-mtls"]
	if mode != CertificateAuthorityMutualTLS {
		t.Errorf("Expected TLSModeMap['ca-mtls'] to be TLSModeCaMTLS, got '%d'", mode)
	}

	mode = TLSModeMap["ss-tls"]
	if mode != SelfSignedTLS {
		t.Errorf("Expected TLSModeMap['ss-tls'] to be TLSModeSsTLS, got '%d'", mode)
	}

	mode = TLSModeMap["ss-mtls"]
	if mode != SelfSignedMutualTLS {
		t.Errorf("Expected TLSModeMap['ss-mtls'] to be TLSModeSsMTLS, got '%d'", mode)
	}

	mode = TLSModeMap["invalid"]
	if mode != TLSMode(0) {
		t.Errorf("Expected TLSModeMap['invalid'] to be TLSMode(0), got '%d'", mode)
	}
}
