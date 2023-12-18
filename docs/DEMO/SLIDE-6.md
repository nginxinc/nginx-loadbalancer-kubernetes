### Configuring TLS in NLK

NLK uses a Kubernetes ConfigMap to configure TLS. The ConfigMap is named `nlk-config` and is located in the `nlk` namespace.

There are three fields in the ConfigMap that are used to configure TLS:
* tls-mode: The TLS mode to use. Valid values are `none`, `ss-tls`, `ss-mtls`, `ca-tls`, and `ca-mtls`.
* ca-certificate: The CA certificate to use. This field is only used when `tls-mode` is set to `ss-tls` or `ss-mtls`. This certificate contains the "Chain of Trust" that is used to verify the authenticity of the TLS certificate.
* client-certificate: The client certificate to use. This field is only used when `tls-mode` is set to `ss-mtls` or `ca-mtls`. This certificate is provided to the NGINX Plus hosts for client authentication.

The fields required depend on the `tls-mode` value.

[Next](SLIDE-7.md)