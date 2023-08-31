# Generate a client certificate

When using self-signed certificates in the `ss-mtls` mode, a certificate needs to be generated for NLK.

The certificate has the same basic field as the CA certificate, with the addition of `clientAuth` in the `extendedKeyUsage` field:

```bash
[ req ]
distinguished_name = dn
prompt = no

[ dn ]
C=[COUNTRY]
ST=[STATE]
L=[LOCALITY]
O=[ORGANIZATION]
OU=[ORGANIZATION_UNIT]

[ client ]
extendedKeyUsage = clientAuth
```

Create a file using the above as a template and update the values in the  `[ dn ]` section; then use following command to generate the certificate request and key:

```bash 
openssl genrsa -out client.key 2048
openssl req -new -key client.key -config client.cnf -out client.csr
```

The output of the above commands will be the client certificate request (`client.csr`) and key (`client.key`).

##### Sign the client certificate

```bash
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365 -sha256 -extfile client.cnf -extensions client
```

The output of the above command will be the client certificate (`client.crt`).

#### Verify the Client Certificate has the correct extendedKeyUsage

```bash
openssl x509 -in client.crt -noout -purpose | grep 'SSL client :'
```

Look for `SSL client : Yes` in the output.

## References

- [Distinguished Name reference](http://certificate.fyicenter.com/2098_OpenSSL_req_-distinguished_name_Configuration_Section.html)