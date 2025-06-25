## TLS in NLK

NLK supports three options for securing communications between the NLK and NGINX Plus:

1. No TLS
2. TLS with self-signed certificates
3. TLS with certificates signed by a Certificate Authority (CA)

Within the TLS options there are two sub-options:

1. One-way TLS
2. Mutual TLS

This gives five options for securing communications between the NLK and NGINX Plus.

* No TLS: No authentication nor encryption is used.
* One-way TLS with self-signed certificates: The NLK verifies the NGINX Plus's identity, but the NGINX Plus does not verify the NLK's identity.
* One-way TLS with certificates signed by a CA: The NLK verifies the NGINX Plus's identity, but the NGINX Plus does not verify the NLK's identity.
* Mutual TLS with self-signed certificates: The NLK verifies the NGINX Plus's identity, and the NGINX Plus verifies the NLK's identity.
* Mutual TLS with certificates signed by a CA: The NLK verifies the NGINX Plus's identity, and the NGINX Plus verifies the NLK's identity.

[Next](SLIDE-6.md)