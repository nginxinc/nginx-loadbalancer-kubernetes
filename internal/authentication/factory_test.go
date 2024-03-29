/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package authentication

import (
	"testing"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

const (
	CaCertificateSecretKey     = "nlk-tls-ca-secret"
	ClientCertificateSecretKey = "nlk-tls-client-secret"
)

func TestTlsFactory_UnspecifiedModeDefaultsToNoTls(t *testing.T) {
	t.Parallel()
	settings := configuration.Settings{}

	tlsConfig, err := NewTLSConfig(&settings)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if tlsConfig == nil {
		t.Fatalf(`tlsConfig should not be nil`)
	}

	if tlsConfig.InsecureSkipVerify != true {
		t.Fatalf(`tlsConfig.InsecureSkipVerify should be true`)
	}
}

func TestTlsFactory_SelfSignedTlsMode(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(caCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.SelfSignedTLS,
		Certificates: &certification.Certificates{
			Certificates:               certificates,
			CaCertificateSecretKey:     CaCertificateSecretKey,
			ClientCertificateSecretKey: ClientCertificateSecretKey,
		},
	}

	tlsConfig, err := NewTLSConfig(&settings)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if tlsConfig == nil {
		t.Fatalf(`tlsConfig should not be nil`)
	}

	if tlsConfig.InsecureSkipVerify != false {
		t.Fatalf(`tlsConfig.InsecureSkipVerify should be false`)
	}

	if len(tlsConfig.Certificates) != 0 {
		t.Fatalf(`tlsConfig.Certificates should be empty`)
	}

	if tlsConfig.RootCAs == nil {
		t.Fatalf(`tlsConfig.RootCAs should not be nil`)
	}
}

func TestTlsFactory_SelfSignedTlsModeCertPoolError(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(invalidCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.SelfSignedTLS,
		Certificates: &certification.Certificates{
			Certificates: certificates,
		},
	}

	_, err := NewTLSConfig(&settings)
	if err == nil {
		t.Fatalf(`Expected an error`)
	}

	if err.Error() != "failed to decode PEM block containing CA certificate" {
		t.Fatalf(`Unexpected error message: %v`, err)
	}
}

func TestTlsFactory_SelfSignedTlsModeCertPoolCertificateParseError(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(invalidCertificateDataPEM())

	settings := configuration.Settings{
		TLSMode: configuration.SelfSignedTLS,
		Certificates: &certification.Certificates{
			Certificates:               certificates,
			CaCertificateSecretKey:     CaCertificateSecretKey,
			ClientCertificateSecretKey: ClientCertificateSecretKey,
		},
	}

	_, err := NewTLSConfig(&settings)
	if err == nil {
		t.Fatalf(`Expected an error`)
	}

	if err.Error() != "error parsing certificate: x509: inner and outer signature algorithm identifiers don't match" {
		t.Fatalf(`Unexpected error message: %v`, err)
	}
}

func TestTlsFactory_SelfSignedMtlsMode(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(caCertificatePEM())
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), clientCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.SelfSignedMutualTLS,
		Certificates: &certification.Certificates{
			Certificates:               certificates,
			CaCertificateSecretKey:     CaCertificateSecretKey,
			ClientCertificateSecretKey: ClientCertificateSecretKey,
		},
	}

	tlsConfig, err := NewTLSConfig(&settings)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if tlsConfig == nil {
		t.Fatalf(`tlsConfig should not be nil`)
	}

	if tlsConfig.InsecureSkipVerify != false {
		t.Fatalf(`tlsConfig.InsecureSkipVerify should be false`)
	}

	if len(tlsConfig.Certificates) == 0 {
		t.Fatalf(`tlsConfig.Certificates should not be empty`)
	}

	if tlsConfig.RootCAs == nil {
		t.Fatalf(`tlsConfig.RootCAs should not be nil`)
	}
}

func TestTlsFactory_SelfSignedMtlsModeCertPoolError(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(invalidCertificatePEM())
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), clientCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.SelfSignedMutualTLS,
		Certificates: &certification.Certificates{
			Certificates: certificates,
		},
	}

	_, err := NewTLSConfig(&settings)
	if err == nil {
		t.Fatalf(`Expected an error`)
	}

	if err.Error() != "failed to decode PEM block containing CA certificate" {
		t.Fatalf(`Unexpected error message: %v`, err)
	}
}

func TestTlsFactory_SelfSignedMtlsModeClientCertificateError(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(caCertificatePEM())
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), invalidCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.SelfSignedMutualTLS,
		Certificates: &certification.Certificates{
			Certificates:               certificates,
			CaCertificateSecretKey:     CaCertificateSecretKey,
			ClientCertificateSecretKey: ClientCertificateSecretKey,
		},
	}

	_, err := NewTLSConfig(&settings)
	if err == nil {
		t.Fatalf(`Expected an error`)
	}

	if err.Error() != "tls: failed to find any PEM data in certificate input" {
		t.Fatalf(`Unexpected error message: %v`, err)
	}
}

func TestTlsFactory_CaTlsMode(t *testing.T) {
	t.Parallel()
	settings := configuration.Settings{
		TLSMode: configuration.CertificateAuthorityTLS,
	}

	tlsConfig, err := NewTLSConfig(&settings)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if tlsConfig == nil {
		t.Fatalf(`tlsConfig should not be nil`)
	}

	if tlsConfig.InsecureSkipVerify != false {
		t.Fatalf(`tlsConfig.InsecureSkipVerify should be false`)
	}

	if len(tlsConfig.Certificates) != 0 {
		t.Fatalf(`tlsConfig.Certificates should be empty`)
	}

	if tlsConfig.RootCAs != nil {
		t.Fatalf(`tlsConfig.RootCAs should be nil`)
	}
}

func TestTlsFactory_CaMtlsMode(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), clientCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.CertificateAuthorityMutualTLS,
		Certificates: &certification.Certificates{
			Certificates:               certificates,
			CaCertificateSecretKey:     CaCertificateSecretKey,
			ClientCertificateSecretKey: ClientCertificateSecretKey,
		},
	}

	tlsConfig, err := NewTLSConfig(&settings)
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if tlsConfig == nil {
		t.Fatalf(`tlsConfig should not be nil`)
	}

	if tlsConfig.InsecureSkipVerify != false {
		t.Fatalf(`tlsConfig.InsecureSkipVerify should be false`)
	}

	if len(tlsConfig.Certificates) == 0 {
		t.Fatalf(`tlsConfig.Certificates should not be empty`)
	}

	if tlsConfig.RootCAs != nil {
		t.Fatalf(`tlsConfig.RootCAs should be nil`)
	}
}

func TestTlsFactory_CaMtlsModeClientCertificateError(t *testing.T) {
	t.Parallel()
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(caCertificatePEM())
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), invalidCertificatePEM())

	settings := configuration.Settings{
		TLSMode: configuration.CertificateAuthorityMutualTLS,
		Certificates: &certification.Certificates{
			Certificates: certificates,
		},
	}

	_, err := NewTLSConfig(&settings)
	if err == nil {
		t.Fatalf(`Expected an error`)
	}

	if err.Error() != "tls: failed to find any PEM data in certificate input" {
		t.Fatalf(`Unexpected error message: %v`, err)
	}
}

// caCertificatePEM returns a PEM-encoded CA certificate.
// Note: The certificate is self-signed and generated explicitly for tests,
// it is not used anywhere else.
func caCertificatePEM() string {
	return `
-----BEGIN CERTIFICATE-----
MIIDTzCCAjcCFA4Zdj3E9TdjOP48eBRDGRLfkj7CMA0GCSqGSIb3DQEBCwUAMGQx
CzAJBgNVBAYTAlVTMRMwEQYDVQQIDApXYXNoaW5ndG9uMRAwDgYDVQQHDAdTZWF0
dGxlMQ4wDAYDVQQKDAVOR0lOWDEeMBwGA1UECwwVQ29tbXVuaXR5ICYgQWxsaWFu
Y2VzMB4XDTIzMDkyOTE3MTY1MVoXDTIzMTAyOTE3MTY1MVowZDELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCldhc2hpbmd0b24xEDAOBgNVBAcMB1NlYXR0bGUxDjAMBgNV
BAoMBU5HSU5YMR4wHAYDVQQLDBVDb21tdW5pdHkgJiBBbGxpYW5jZXMwggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCwlI4ZvJ/6hvqULFVL+1ZSRDTPQ48P
umehJhPz6xPhC9UkeTe2FZxm2Rsi1I5QXm/bTG2OcX775jgXzae9NQjctxwrz4Ks
LOWUvRkkfhQR67xk0Noux76/9GWGnB+Fapn54tlWql6uHQfOu1y7MCRkZ27zHbkk
lq4Oa2RmX8rIyECWgbTyL0kETBVJU8bYORQ5JjhRlz08inq3PggY8blrehIetrWN
dw+gzcqdvAI2uSCodHTHM/77KipnYmPiSiDjSDRlXdxTG8JnyIB78IoH/sw6RyBm
CvVa3ytvKziXAvbBoXq5On5WmMRF97p/MmBc53ExMuDZjA4fisnViS0PAgMBAAEw
DQYJKoZIhvcNAQELBQADggEBAJeoa2P59zopLjBInx/DnWn1N1CmFLb0ejKxG2jh
cOw15Sx40O0XrtrAto38iu4R/bkBeNCSUILlT+A3uYDila92Dayvls58WyIT3meD
G6+Sx/QDF69+4AXpVy9mQ+hxcofpFA32+GOMXwmk2OrAcdSkkGSBhZXgvTpQ64dl
xSiQ5EQW/K8LoBoEOXfjIZJNPORgKn5MI09AY7/47ycKDKTUU2yO8AtIHYKttw0x
kfIg7QOdo1F9IXVpGjJI7ynyrgsCEYxMoDyH42Dq84eKgrUFLEXemEz8hgdFgK41
0eUYhAtzWHbRPBp+U/34CQoZ5ChNFp2YipvtXrzKE8KLkuM=
-----END CERTIFICATE-----
`
}

func invalidCertificatePEM() string {
	return `
-----BEGIN CERTIFICATE-----
MIIClzCCAX+gAwIBAgIJAIfPhC0RG6CwMA0GCSqGSIb3DQEBCwUAMBkxFzAVBgNV
BAMMDm9pbCBhdXRob3JpdHkwHhcNMjAwNDA3MTUwOTU1WhcNMjEwNDA2MTUwOTU1
WjBMMSAwHgYDVQQLDBd5b3VuZy1jaGFsbGVuZ2UgdGVzdCBjb25zdW1lczEfMB0G
A1UECwwWc28wMS5jb3Jwb3JhdGlvbnNvY2lhbDEhMB8GA1UEAwwYc29tMS5jb3Jw
b3JhdGlvbnNvY2lhbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQDGRX31uzy+yLUOz7wOJHHm2dzrDgUbC6RZDjURvZxyt2Zi5wYWsEB5r5YhN7L0
y1R9f+MGwNITIz9nYZuU/PLFOvzF5qX7A8TbdgjZEqvXe2NZ9J2z3iWvYQLN8Py3
nv/Y6wadgXEBRCNNuIg/bQ9XuOr9tfB6j4Ut1GLU0eIlV/L3Rf9Y6SgrAl+58ITj
Wrg3Js/Wz3J2JU4qBD8U4I3XvUyfnX2SAG8Llm4KBuYz7g63Iu05s6RnmG+Xhu2T
5f2DWZUeATWbAlUW/M4NLO1+5H0gOr0TGulETQ6uElMchT7s/H6Rv1CV+CNCCgEI
adRjWJq9yQ+KrE+urSMCXu8XAgMBAAGjUzBRMB0GA1UdDgQWBBRb40pKGU4lNvqB
1f5Mz3t0N/K3hzAfBgNVHSMEGDAWgBRb40pKGU4lNvqB1f5Mz3t0N/K3hzAPBgNV
HREECDAGhwQAAAAAAAAwCgYIKoZIzj0EAwIDSAAwRQIhAP3ST/mXyRXsU2ciRoE
gE6trllODFY+9FgT6UbF2TwzAiAAuaUxtbk6uXLqtD5NtXqOQf0Ckg8GQxc5V1G2
9PqTXQ==
-----END CERTIFICATE-----
`
}

// Yoinked from https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/crypto/x509/x509_test.go, line 3385
// This allows the `buildCaCertificatePool(...)` --> `x509.ParseCertificate(...)` call error branch to be covered.
func invalidCertificateDataPEM() string {
	return `
-----BEGIN CERTIFICATE-----
MIIBBzCBrqADAgECAgEAMAoGCCqGSM49BAMCMAAwIhgPMDAwMTAxMDEwMDAwMDBa
GA8wMDAxMDEwMTAwMDAwMFowADBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABOqV
EDuVXxwZgIU3+dOwv1SsMu0xuV48hf7xmK8n7sAMYgllB+96DnPqBeboJj4snYnx
0AcE0PDVQ1l4Z3YXsQWjFTATMBEGA1UdEQEB/wQHMAWCA2FzZDAKBggqhkjOPQQD
AwNIADBFAiBi1jz/T2HT5nAfrD7zsgR+68qh7Erc6Q4qlxYBOgKG4QIhAOtjIn+Q
tA+bq+55P3ntxTOVRq0nv1mwnkjwt9cQR9Fn
-----END CERTIFICATE-----
`
}

// clientCertificatePEM returns a PEM-encoded client certificate.
// Note: The certificate is self-signed and generated explicitly for tests,
// it is not used anywhere else.
func clientCertificatePEM() string {
	return `
-----BEGIN CERTIFICATE-----
MIIEDDCCAvSgAwIBAgIULDFXwGrTohN/PRao2rSLk9VxFdgwDQYJKoZIhvcNAQEL
BQAwXTELMAkGA1UEBhMCVVMxEzARBgNVBAgMCldhc2hpbmd0b24xEjAQBgNVBAcM
CUluZGlhbm9sYTEPMA0GA1UECgwGV2FnbmVyMRQwEgYDVQQLDAtEZXZlbG9wbWVu
dDAeFw0yMzA5MjkxNzA3NTRaFw0yNDA5MjgxNzA3NTRaMGQxCzAJBgNVBAYTAlVT
MRMwEQYDVQQIDApXYXNoaW5ndG9uMRAwDgYDVQQHDAdTZWF0dGxlMQ4wDAYDVQQK
DAVOR0lOWDEeMBwGA1UECwwVQ29tbXVuaXR5ICYgQWxsaWFuY2VzMIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoqNuEZ6+TcFrmzcwp8u8mzk0jPd47GKk
H9wwdkFCzGdd8KJkFQhzLyimZIWkRDYmhaxZd76jKGBpdfyivR4e4Mi5WYlpPGMI
ppM7/rMYP8yn04tkokAazbqjOTlF8NUKqGQwqAN4Z/PvoG2HyP9omGpuLWTbjKto
oGr5aPBIhzlICU3OjHn6eKaekJeAYBo3uQFYOxCjtE9hJLDOY4q7zomMJfYoeoA2
Afwkx1Lmozp2j/esB52/HlCKVhAOzZsPzM+E9eb1Q722dUed4OuiVYSfrDzeImrA
TufzTBTMEpFHCtdBGocZ3LRd9qmcP36ZCMsJNbYnQZV3XsI4JhjjHwIDAQABo4G8
MIG5MBMGA1UdJQQMMAoGCCsGAQUFBwMCMB0GA1UdDgQWBBRDl4jeiE1mJDPrYmQx
g2ndkWxpYjCBggYDVR0jBHsweaFhpF8wXTELMAkGA1UEBhMCVVMxEzARBgNVBAgM
Cldhc2hpbmd0b24xEjAQBgNVBAcMCUluZGlhbm9sYTEPMA0GA1UECgwGV2FnbmVy
MRQwEgYDVQQLDAtEZXZlbG9wbWVudIIUNxx2Mr+PKXiF3d2i51fb/rnWbBgwDQYJ
KoZIhvcNAQELBQADggEBAL0wS6LkFuqGDlhaTGnAXRwRDlC6uwrm8wNWppaw9Vqt
eaZGFzodcCFp9v8jjm1LsTv7gEUBnWtn27LGP4GJSpZjiq6ulJypBxo/G0OkMByK
ky4LeGY7/BQzjzHdfXEq4gwfC45ni4n54uS9uzW3x+AwLSkxPtBxSwxhtwBLo9aE
Ql4rHUoWc81mhGO5mMZBaorxZXps1f3skfP+wZX943FIMt5gz4hkxwFp3bI/FrqH
R8DLUlCzBA9+7WIFD1wi25TV+Oyq3AjT/KiVmR+umrukhnofCWe8JiVpb5iJcd2k
Rc7+bvyb5OCnJdEX08XGWmF2/OFKLrCzLH1tQxk7VNE=
-----END CERTIFICATE-----
`
}

// clientKeyPEM returns a PEM-encoded client key.
// Note: The key is self-signed and generated explicitly for tests,
// it is not used anywhere else.
func clientKeyPEM() string {
	return `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCio24Rnr5NwWub
NzCny7ybOTSM93jsYqQf3DB2QULMZ13womQVCHMvKKZkhaRENiaFrFl3vqMoYGl1
/KK9Hh7gyLlZiWk8Ywimkzv+sxg/zKfTi2SiQBrNuqM5OUXw1QqoZDCoA3hn8++g
bYfI/2iYam4tZNuMq2igavlo8EiHOUgJTc6Mefp4pp6Ql4BgGje5AVg7EKO0T2Ek
sM5jirvOiYwl9ih6gDYB/CTHUuajOnaP96wHnb8eUIpWEA7Nmw/Mz4T15vVDvbZ1
R53g66JVhJ+sPN4iasBO5/NMFMwSkUcK10EahxnctF32qZw/fpkIywk1tidBlXde
wjgmGOMfAgMBAAECggEAA+R2b2yFsHW3HhVhkDqDjpF9bPxFRB8OP4b1D/d64kp9
CJPSYmB75T6LUO+T4WAMZvmbgI6q9/3quDyuJmmQop+bNAXiY2QZYmc2sd9Wbrx2
rczxwSJYoeDcJDP3NQ7cPPB866B9ortHWmcUr15RgghWD7cQvBqkG+bDhlvt2HKg
NZmL6R0U1bVAlRMtFJiEdMHuGnPmoDU5IGc1fKjsgijLeMboUrEaXWINoEm8ii5e
/mnsfLCBmeJAsKuXxL8/1UmvWYE/ltDfYBVclKhcH2UWTZv7pdRtHnu49lkZivUB
ZvH2DHsSMjXj6+HHr6RcRGmnMDyfhJFPCjOdTjf4oQKBgQDeYLWZx22zGXgfb7md
MhdKed9GxMJHzs4jDouqrHy0w95vwMi7RXgeKpKXiCruqSEB/Trtq01f7ekh0mvJ
Ys0h4A5tkrT5BVVBs+65uF/kSF2z/CYGNRhAABO7UM+B1e3tlnjfjeb/M78IcFbT
FyBN90A/+a9JGZ4obt3ack3afwKBgQC7OncnXC9L5QCWForJWQCNO3q3OW1Gaoxe
OAnmnPSJ7NUd7xzDNE8pzBUWXysZCoRU3QNElcQfzHWtZx1iqJPk3ERK2awNsnV7
X2Fu4vHzIr5ZqVnM8NG7+iWrxRLf+ctcEvPiqRYo+g+r5tTGJqWh2nh9W7iQwwwE
1ikoxFBnYQKBgCbDdOR5fwXZSrcwIorkUGsLE4Cii7s4sXYq8u2tY4+fFQcl89ex
JF8dzK/dbJ5tnPNb0Qnc8n/mWN0scN2J+3gMNnejOyitZU8urk5xdUW115+oNHig
iLmfSdE9JO7c+7yOnkNZ2QpjWsl9y6TAQ0FT+D8upv93F7q0mLebdTbBAoGBALmp
r5EThD9RlvQ+5F/oZ3imO/nH88n5TLr9/St4B7NibLAjdrVIgRwkqeCmfRl26WUy
SdRQY81YtnU/JM+59fbkSsCi/FAU4RV3ryoD2QRPNs249zkYshMjawncAuyiS/xB
OyJQpI3782B3JhZdKrDG8eb19p9vG9MMAILRsh3hAoGASCvmq10nHHGFYTerIllQ
sohNaw3KDlQTkpyOAztS4jOXwvppMXbYuCznuJbHz0NEM2ww+SiA1RTvD/gosYYC
mMgqRga/Qu3b149M3wigDjK+RAcyuNGZN98bqU/UjJLjqH6IMutt59+9XNspcD96
z/3KkMx4uqJXZyvQrmkolSg=
-----END PRIVATE KEY-----
`
}

func buildClientCertificateEntry(keyPEM, certificatePEM string) map[string]core.SecretBytes {
	return map[string]core.SecretBytes{
		certification.CertificateKey:    core.SecretBytes([]byte(certificatePEM)),
		certification.CertificateKeyKey: core.SecretBytes([]byte(keyPEM)),
	}
}

func buildCaCertificateEntry(certificatePEM string) map[string]core.SecretBytes {
	return map[string]core.SecretBytes{
		certification.CertificateKey: core.SecretBytes([]byte(certificatePEM)),
	}
}
