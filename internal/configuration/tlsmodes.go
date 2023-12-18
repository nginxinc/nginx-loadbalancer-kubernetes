/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package configuration

const (
	NoTLS TLSMode = iota
	CertificateAuthorityTLS
	CertificateAuthorityMutualTLS
	SelfSignedTLS
	SelfSignedMutualTLS
)

const (
	NoTLSString                         = "no-tls"
	CertificateAuthorityTLSString       = "ca-tls"
	CertificateAuthorityMutualTLSString = "ca-mtls"
	SelfSignedTLSString                 = "ss-tls"
	SelfSignedMutualTLSString           = "ss-mtls"
)

type TLSMode int

var TLSModeMap = map[string]TLSMode{
	NoTLSString:                         NoTLS,
	CertificateAuthorityTLSString:       CertificateAuthorityTLS,
	CertificateAuthorityMutualTLSString: CertificateAuthorityMutualTLS,
	SelfSignedTLSString:                 SelfSignedTLS,
	SelfSignedMutualTLSString:           SelfSignedMutualTLS,
}

func (t TLSMode) String() string {
	modes := []string{
		NoTLSString,
		CertificateAuthorityTLSString,
		CertificateAuthorityMutualTLSString,
		SelfSignedTLSString,
		SelfSignedMutualTLSString,
	}
	if t < NoTLS || t > SelfSignedMutualTLS {
		return ""
	}
	return modes[t]
}
