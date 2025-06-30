/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package configuration

const (
	// NoTLS is deprecated as misleading. It is the same as SkipVerifyTLS.
	NoTLS = "no-tls"
	// SkipVerifyTLS causes the http client to skip verification of the NGINX
	// host's certificate chain and host name.
	SkipVerifyTLS = "skip-verify-tls"
	// CertificateAuthorityTLS is deprecated as misleading. This is the same as
	// the default behavior which is to verify the NGINX hosts's certificate
	// chain and host name, if https is used.
	CertificateAuthorityTLS = "ca-tls"
)

var tlsModeMap = map[string]bool{
	NoTLS:                   true,
	SkipVerifyTLS:           true,
	CertificateAuthorityTLS: false,
}
