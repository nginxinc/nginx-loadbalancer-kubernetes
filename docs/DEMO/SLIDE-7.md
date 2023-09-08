### How NLK uses TLS

NLK uses the `github.com/nginxinc/nginx-plus-go-client` to communicate with the NGINX Plus API. This library accepts a low-level HTTP client that is used for the actual communication with the NGINX Plus hosts. NLK generates a TLS Configuration based on the configured `tls-mode` and provides it to the low-level HTTP client.

[Next](SLIDE-8.md)