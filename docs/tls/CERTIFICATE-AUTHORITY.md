# Generate a Certificate Authority (CA) 

When using self-signed certificates, the first step is to generate the Certificate Authority (CA).

The following commands will generate the CA certificate and key:

```bash
openssl req -newkey rsa:2048 -nodes -x509 -out ca.crt -keyout ca.key
```

You will be prompted to enter the Distinguished Name (DN) information for the CA.

Alternatively, you can provide the DN information in a file, an example is shown below:

```bash
[ req ]
distinguished_name = dn
prompt = no
req_extensions = req_ext

[ req_ext ]
basicConstraints = CA:TRUE
keyUsage = critical, keyCertSign, cRLSign

[ dn ]
C=[COUNTRY]
ST=[STATE]
L=[LOCALITY]
O=[ORGANIZATION]
OU=[ORGANIZATION_UNIT]
```

Create a file using the above as a template and update the values in the  `[ dn ]` section; then use following command to generate the CA certificate and key:

```bash
openssl req -newkey rsa:2048 -nodes -x509 -config ca.cnf -out ca.crt -keyout ca.key
```

The output of the above command will be the CA certificate (`ca.crt`) and key (`ca.key`).

## References

- [Distinguished Name reference](http://certificate.fyicenter.com/2098_OpenSSL_req_-distinguished_name_Configuration_Section.html)

