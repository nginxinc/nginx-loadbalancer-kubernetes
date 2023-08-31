# Configuring NGINX Plus

## Configuring the server certificate

Depending on the mode chosen, some files will need to be copied to the NGINX Plus host. The following table shows which files are required for each mode:

| Mode | Server Certificate | CA Certificate     |
| ---- | ------------------ |--------------------|
| no-tls | |                    | |
| ss-tls | :heavy_check_mark: | :heavy_check_mark: |
| ss-mtls | :heavy_check_mark: | :heavy_check_mark: |
| ca-tls | :heavy_check_mark: |                    |
| ca-mtls | :heavy_check_mark: | |


Copy the necessary server files for the chosen mode to the NGINX host; place them in the `/etc/ssl/certs/nginx` directory.

### One-way TLS

The following applies to all modes other than `no-tls`. 

To configure NGINX Plus to use the `server.crt` and `server.key` files for TLS,  
add the following to the `http` or `server` context in the `/etc/nginx/nginx.conf` file:

```bash
http {
  ssl_certificate       /etc/ssl/certs/nginx/server.crt;
  ssl_certificate_key   /etc/ssl/certs/nginx/server.key;
}
```

For more information about the `ssl_certificate` directive, refer to the NGINX [documentation](https://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_certificate).

Reload the NGINX Plus configuration to apply the changes.

```bash
nginx -s reload
```

### Mutual TLS

In the `/etc/nginx/nginx.conf` file, add the following to the `http` or `server` context:

#### Self-signed certificates

When using `ss-mtls` mode, the CA certificate must be provided to NGINX Plus:

```bash
http {
  ssl_client_certificate    /etc/ssl/certs/nginx/ca.crt;
  ssl_verify_client         on;
  ssl_verify_depth          3;
}
```

This will configure NGINX Plus to use the `ca.crt` file for client authentication.

#### CA-signed certificates

When using `ca-mtls` mode, the `ssl_client_certificate` directive is not required:

```bash
http {
  ssl_verify_client         on;
  ssl_verify_depth          3;
}
```

Please refer to the [NGINX documentation](https://nginx.org/en/docs/http/ngx_http_ssl_module.html#ssl_client_certificate) for details on the `ssl_client_certificate` directive.

Reload the NGINX Plus configuration to apply the changes.

```bash
nginx -s reload
```

Test with curl:

```bash
curl --cert client.crt --key client.key --cacert ca.crt https://<your-host>/api
```