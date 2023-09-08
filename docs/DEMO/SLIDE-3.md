## TLS uses Certificates

A Transport Layer Security (TLS) certificate, also known as an SSL certificate (Secure Sockets Layer), is a digital document that plays a crucial role in securing internet communications. Imagine it as a special, electronic passport for websites.

A TLS certificate includes a "Chain of Trust" where there are multiple certificates that are used to verify the authenticity of the certificate. The certificate at the top of the chain is called the Root Certificate Authority (CA) Certificate. The Root CA Certificate is used to sign the certificates below it in the chain. The Root CA Certificate is not signed by any other certificate in the chain.

The Root CA Certificate is used to sign the Intermediate CA Certificate. The Intermediate CA Certificate is used to sign the TLS Certificate. The TLS Certificate is used to sign the TLS Certificate Signing Request (CSR). The TLS Certificate Signing Request is used to sign the TLS Certificate.

Root CA certificates can be expensive and are not required for most use cases. An alternative to purchasing an Intermediate CA certificate is to use a self-signed certificate. Self-signed certificates are free and can be used to sign TLS certificates.

[Next](SLIDE-4.md)