# Certificate Authority Mutual TLS Mode

This mode allows NLK to verify it is connecting to the correct NGINX Plus instance, allows NGINX Plus to verify it is connecting to the correct NLK, and encrypts the data between NLK and NGINX Plus.

## Overview

Mutual TLS is used to encrypt the traffic between NLK and NGINX Plus, to ensure NLK verifies the NGINX Plus server, and to ensure NGINX Plus verifies NLK.

## Certificates

To configure this mode, the following certificates are required:

- Server Certificate
- Client Certificate

See the following sections for instructions on how to create these certificates.

### Certificate Authority (CA)

Provided by the user.

### Server Certificate (NGINX Plus)

Use your own certificate authority (CA) to generate a server certificate and key.

### Client Certificate (NLK)

Use your own certificate authority (CA) to generate a client certificate and key.

## Kubernetes Secrets

NLK accesses the necessary certificates for each mode from Kubernetes Secrets. For this mode, the following Kubernetes Secret(s) are required:
- Client Certificate

To create the Kubernetes Secret containing the CA certificate, refer to the [Kubernetes Secrets](./KUBERNETES-SECRETS.md) guide;
the name and location of the certificate(s) created above should be used. The name of the Secret will be needed for the ConfigMap (discussed below).

## ConfigMap


NLK is configured via a ConfigMap. The ConfigMap is named `nlk-config` and is located in the `nlk` namespace.

Depending on which mode is chosen, certain fields will need to be updated in the NLK ConfigMap.

For this mode, the `mode` and `clientCertificate` fields need to be included. The `mode` field should be set to `ca-mtls` 
and the `clientCertificate` field should be set to the name of the Kubernetes Secret containing the Client certificate created above.

The following is an example of a ConfigMap for this mode (be sure to update the `nginx-hosts` field with the correct NGINX Plus API endpoints):

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nlk-config
  namespace: nlk
data:
  nginx-hosts: "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api"
  mode: "ca-mtls"
  clientCertificate: "nlk-tls-client-secret"
```

## Deployment

Save the above ConfigMap definition to a file named `ca-mtls-configmap.yaml`, then deploy the ConfigMap using the following command:

```bash
kubectl apply -f docs/tls/ca-mtls-configmap.yaml
```

## Configuring NGINX Plus

Refer to the [NGINX Plus Configuration](./NGINX-PLUS-CONFIGURATION.md) guide for instructions on configuring NGINX Plus to use the certificates created above.

The steps in both the ["One-way TLS"](./NGINX-PLUS-CONFIGURATION.md#one-way-tls) section and the ["Mutual TLS"](./NGINX-PLUS-CONFIGURATION.md#mutual-tls) section are required for this mode.

## Verification

To verify the ConfigMap was deployed correctly, run the following command:

```bash
kubectl get configmap -n nlk nlk-config -o yaml
```

The output should match the ConfigMap above.

To verify NLK is running, follow the instructions in either the [TCP](../tcp/tcp-installation-guide.md) or [HTTP](../http/http-installation-guide.md) guides.

