# One-way TLS with self-signed certificates

This mode allows NLK to verify it is connecting to the correct NGINX Plus instance, and encrypts the data between NLK and NGINX Plus.

## Overview

One-way TLS is used to encrypt the traffic between NLK and NGINX Plus, and to ensure NLK verifies the NGINX Plus server;
however, the NGINX Plus server _does not_ validate NLK.

## Certificates

To configure this mode, the following certificates are required:

- CA Certificate
- Server Certificate

See the following sections for instructions on how to create these certificates.

### Certificate Authority (CA)

Follow the instructions in the [Certificate Authority](./CERTIFICATE-AUTHORITY.md) guide to create a self-signed CA certificate and key.

### Server Certificate (NGINX Plus)

Follow the instructions in the [Server Certificate](./SERVER-CERTIFICATE.md) guide to create a self-signed server certificate and key.

## Kubernetes Secrets

NLK accesses the necessary certificates for each mode from Kubernetes Secrets. For this mode, the following Kubernetes Secret(s) are required:
- CA Certificate

To create the Kubernetes Secret containing the CA certificate, refer to the [Kubernetes Secrets](./KUBERNETES-SECRETS.md) guide;
the name and location of the certificate(s) created above should be used. The name of the Secret will be needed for the ConfigMap (discussed below).

## ConfigMap

NLK is configured via a ConfigMap. The ConfigMap is named `nlk-config` and is located in the `nlk` namespace.

Depending on which mode is chosen, certain fields will need to be updated in the NLK ConfigMap.

For this mode, both the `tlsMode` and `caCertificates` fields need to be included. The `tlsMode` field should be set to `ss-tls`, 
and the `caCertificates` field should be set to the name of the Kubernetes Secret containing the CA certificate created above.

The following is an example of a ConfigMap for this mode (be sure to update the `nginx-hosts` field with the correct NGINX Plus API endpoints):


```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nlk-config
  namespace: nlk
data:
  nginx-hosts: "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api"
  tlsMode: "ss-tls"
  caCertificate: "nlk-tls-ca-secret"
```

## Deployment

Save the above ConfigMap definition to a file named `ss-tls-configmap.yaml`, then deploy the ConfigMap using the following command:

```bash
kubectl apply -f docs/tls/ss-tls-configmap.yaml
```

## Configuring NGINX Plus

Refer to the [NGINX Plus Configuration](./NGINX-PLUS-CONFIGURATION.md) guide for instructions on configuring NGINX Plus to use the certificates created above.

Only the steps in the ["One-way TLS"](./NGINX-PLUS-CONFIGURATION.md#one-way-tls) section are required for this mode.

## Verification

To verify the ConfigMap was deployed correctly, run the following command:

```bash
kubectl get configmap -n nlk nlk-config -o yaml
```

The output should match the ConfigMap above.

To verify NLK is running, follow the instructions in either the [TCP](../tcp/tcp-installation-guide.md) or [HTTP](../http/http-installation-guide.md) guides.
