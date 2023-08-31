# Generate a server certificate

When using self-signed certificates in the `ss-tls` and `ss-mtls` modes, a certificate needs to be generated for the NGINX Plus server.

The certificate has the same basic fields as the CA certificate, but with a few additional fields.

The following is an example of a configuration file for server certificates:

```bash
[ req ]
distinguished_name = dn
req_extensions = req_ext
prompt = no

[ dn ]
C=[COUNTRY]
ST=[STATE]
L=[LOCALITY]
O=[ORGANIZATION]
OU=[ORGANIZATION_UNIT]
CN=[COMMON_NAME]

[ req_ext ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names
extendedKeyUsage = serverAuth

[ alt_names ]
DNS.1 = mydomain.com
DNS.2 = server.mydomain.com
DNS.3 = *.mydomain.com
IP.1 = 10.0.0.10
IP.2 = 10.0.0.11
```

This example includes sensible defaults in terms of extensions, to learn more see the [OpenSSL extensions documentation](https://www.openssl.org/docs/manmaster/man5/x509v3_config.html).

Create a file using this as a template, the following example commands use the name `server.cnf`.

Be sure to update the Distinguished Name (DN) information and the SAN information (DNS / IP entries in the `alt_names` section) as appropriate.
Doing so ensures that the certificate is valid for the server by providing an unambiguous match between the server (IP addresses and/or domain names) and the certificate.
A reference for the DN fields can be found [here](http://certificate.fyicenter.com/2098_OpenSSL_req_-distinguished_name_Configuration_Section.html).


The following commands will generate the server certificate and key:

```bash
openssl genrsa -out server.key 2048
openssl req -new -key server.key -config server.cnf -out server.csr
```

The output of the above commands will be the server certificate request (`server.csr`) and key (`server.key`).

##### Sign the server certificate

```bash
openssl x509  -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extensions req_ext -extfile server.cnf
```

The output of the above command will be the server certificate (`server.crt`).

##### Verify the Server Certificate has the SAN

```bash
openssl x509 -in server.crt -text -noout | grep -A 1 "Subject Alternative Name"
```

Look for the DNS / IP entries in the `Subject Alternative Name` section in the output; the values entered in the `alt_names` section of the `server.cnf` file should be present.
