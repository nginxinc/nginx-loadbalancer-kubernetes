/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package certification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

const (
	CaCertificateSecretKey = "nlk-tls-ca-secret"
)

func TestNewCertificate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	certificates := NewCertificates(ctx, nil)

	if certificates == nil {
		t.Fatalf(`certificates should not be nil`)
	}
}

func TestCertificates_Initialize(t *testing.T) {
	t.Parallel()
	certificates := NewCertificates(context.Background(), nil)

	err := certificates.Initialize()
	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}
}

func TestCertificates_RunWithoutInitialize(t *testing.T) {
	t.Parallel()
	certificates := NewCertificates(context.Background(), nil)

	err := certificates.Run()
	if err == nil {
		t.Fatalf(`Expected error`)
	}

	if err.Error() != `initialize must be called before Run` {
		t.Fatalf(`Unexpected error: %v`, err)
	}
}

func TestCertificates_EmptyCertificates(t *testing.T) {
	t.Parallel()
	certificates := NewCertificates(context.Background(), nil)

	err := certificates.Initialize()
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

func TestCertificates_ExerciseHandlers(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	k8sClient := fake.NewSimpleClientset()

	certificates := NewCertificates(ctx, k8sClient)

	_ = certificates.Initialize()

	certificates.CaCertificateSecretKey = CaCertificateSecretKey

	//nolint:govet,staticcheck
	go func() {
		err := certificates.Run()
		assert.NoError(t, err, "expected no error running certificates")
	}()

	cache.WaitForCacheSync(ctx.Done(), certificates.informer.HasSynced)

	secret := buildSecret()

	/* -- Test Create -- */

	created, err := k8sClient.CoreV1().Secrets(SecretsNamespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf(`error creating the Secret: %v`, err)
	}

	if created.Name != secret.Name {
		t.Fatalf(`Expected name %v, got %v`, secret.Name, created.Name)
	}

	time.Sleep(2 * time.Second)

	caBytes := certificates.GetCACertificate()
	if caBytes == nil {
		t.Fatalf(`Expected non-nil CA certificate`)
	}

	/* -- Test Update -- */

	secret.Labels = map[string]string{"updated": "true"}
	_, err = k8sClient.CoreV1().Secrets(SecretsNamespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf(`error updating the Secret: %v`, err)
	}

	time.Sleep(2 * time.Second)

	caBytes = certificates.GetCACertificate()
	if caBytes == nil {
		t.Fatalf(`Expected non-nil CA certificate`)
	}

	/* -- Test Delete -- */

	err = k8sClient.CoreV1().Secrets(SecretsNamespace).Delete(ctx, secret.Name, metav1.DeleteOptions{})
	if err != nil {
		t.Fatalf(`error deleting the Secret: %v`, err)
	}

	time.Sleep(2 * time.Second)

	caBytes = certificates.GetCACertificate()
	if caBytes != nil {
		t.Fatalf(`Expected nil CA certificate, got: %v`, caBytes)
	}
}

func buildSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CaCertificateSecretKey,
			Namespace: SecretsNamespace,
		},
		Data: map[string][]byte{
			CertificateKey:    []byte(certificatePEM()),
			CertificateKeyKey: []byte(keyPEM()),
		},
		Type: corev1.SecretTypeTLS,
	}
}

// certificatePEM returns a PEM-encoded client certificate.
// Note: The certificate is self-signed and generated explicitly for tests,
// it is not used anywhere else.
func certificatePEM() string {
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

// keyPEM returns a PEM-encoded client key.
// Note: The key is self-signed and generated explicitly for tests,
// it is not used anywhere else.
func keyPEM() string {
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
