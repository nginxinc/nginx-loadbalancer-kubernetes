package main

import (
	"bufio"
	"log/slog"
	"os"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/authentication"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

const (
	CaCertificateSecretKey     = "nlk-tls-ca-secret"
	ClientCertificateSecretKey = "nlk-tls-client-secret"
)

type TLSConfiguration struct {
	Description string
	Settings    configuration.Settings
}

func main() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	configurations := buildConfigMap()

	for name, settings := range configurations {
		slog.Info("\033[H\033[2J")

		slog.Info("\n\n\t*** Building TLS config\n\n", "name", name)

		tlsConfig, err := authentication.NewTLSConfig(settings.Settings)
		if err != nil {
			panic(err)
		}

		rootCaCount := 0
		certificateCount := 0

		if tlsConfig.RootCAs != nil {
			rootCaCount = len(tlsConfig.RootCAs.Subjects()) //nolint:staticcheck
		}

		if tlsConfig.Certificates != nil {
			certificateCount = len(tlsConfig.Certificates)
		}

		slog.Info("Successfully built TLS config",
			"description", settings.Description,
			"rootCaCount", rootCaCount,
			"certificateCount", certificateCount,
		)

		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	slog.Info("\033[H\033[2J")
	slog.Info("\n\n\t*** All done! ***\n\n")
}

func buildConfigMap() map[string]TLSConfiguration {
	configurations := make(map[string]TLSConfiguration)

	configurations["ss-tls"] = TLSConfiguration{
		Description: "Self-signed TLS requires just a CA certificate",
		Settings:    ssTLSConfig(),
	}

	configurations["ss-mtls"] = TLSConfiguration{
		Description: "Self-signed mTLS requires a CA certificate and a client certificate",
		Settings:    ssMtlsConfig(),
	}

	configurations["ca-tls"] = TLSConfiguration{
		Description: "CA TLS requires no certificates",
		Settings:    caTLSConfig(),
	}

	configurations["ca-mtls"] = TLSConfiguration{
		Description: "CA mTLS requires a client certificate",
		Settings:    caMtlsConfig(),
	}

	return configurations
}

func ssTLSConfig() configuration.Settings {
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(caCertificatePEM())
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), clientCertificatePEM())

	return configuration.Settings{
		TLSMode: configuration.SelfSignedTLS,
		Certificates: &certification.Certificates{
			Certificates: certificates,
		},
	}
}

func ssMtlsConfig() configuration.Settings {
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[CaCertificateSecretKey] = buildCaCertificateEntry(caCertificatePEM())
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), clientCertificatePEM())

	return configuration.Settings{
		TLSMode: configuration.SelfSignedMutualTLS,
		Certificates: &certification.Certificates{
			Certificates: certificates,
		},
	}
}

func caTLSConfig() configuration.Settings {
	return configuration.Settings{
		TLSMode: configuration.CertificateAuthorityTLS,
	}
}

func caMtlsConfig() configuration.Settings {
	certificates := make(map[string]map[string]core.SecretBytes)
	certificates[ClientCertificateSecretKey] = buildClientCertificateEntry(clientKeyPEM(), clientCertificatePEM())

	return configuration.Settings{
		TLSMode: configuration.CertificateAuthorityMutualTLS,
		Certificates: &certification.Certificates{
			Certificates: certificates,
		},
	}
}

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
