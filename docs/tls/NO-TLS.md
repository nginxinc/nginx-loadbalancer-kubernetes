# No TLS Mode

This mode is the easiest to configure, but provides no security. Choose this for development environments, or test 
environments where security is not strictly necessary. 

Note: you should test with one of the TLS / mTLS modes before deploying to production to ensure familiarity with the configuration and process.

## Overview

This is the default mode of operation for NLK. It offers no verification of either side of the connection, nor any encryption of the data.

## Certificates

No certificates are required for this mode.

## Kubernetes Secrets

No Kubernetes Secrets are required for this mode.

## ConfigMap

NLK is configured via a ConfigMap. The ConfigMap is named `nlk-config` and is located in the `nlk` namespace. Depending on which mode is chosen, certain fields will need to be updated in the NLK ConfigMap. 

For this mode, only the `tlsMode` field needs to be included, and should be set to `no-tls` (or omitted altogether as this is the default mode).

The following is an example of a ConfigMap for this mode (be sure to update the `nginx-hosts` field with the correct NGINX Plus API endpoints)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nlk-config
  namespace: nlk
data:
  nginx-hosts: "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api"
  tlsMode: "no-tls"
```

## Deployment

Save the above ConfigMap definition to a file named `no-tls-configmap.yaml`, then deploy the ConfigMap using the following command:

```bash
kubectl apply -f docs/tls/no-tls-configmap.yaml
```

## Verification

To verify the ConfigMap was deployed correctly, run the following command:

```bash
kubectl get configmap -n nlk nlk-config -o yaml
```

The output should match the ConfigMap above.

To verify NLK is running, follow the instructions in either the [TCP](../tcp/tcp-installation-guide.md) or [HTTP](../http/http-installation-guide.md) guides.
