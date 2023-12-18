# Kubernetes Secrets

## Overview

Kubernetes Secrets are used to provide the required certificates to NLK. There are two ways to create the Secrets:
- Using `kubectl`
- Using yaml files

The filenames for the certificates created are required for both methods. The examples below assume the certificates were generated in `/tmp` 
and follow the naming conventions in the documentation. 

## Using `kubectl`

The easiest way to create the Secret(s) is by using `kubectl`:

```bash
kubectl create secret tls -n nlk nlk-tls-ca-secret --cert=/tmp/ca.crt --key=/tmp/ca.key
kubectl create secret tls -n nlk nlk-tls-server-secret --cert=/tmp/server.crt --key=/tmp/server.key
kubectl create secret tls -n nlk nlk-tls-client-secret --cert=/tmp/client.crt --key=/tmp/client.key
```

## Using yaml files

The Secrets can also be created using yaml files. The following is an example of a yaml file for the Client Secret (note that the `data` values are truncated):

```yaml
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVCVEN...
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0l...
kind: Secret
metadata:
  name: nlk-tls-ca-secret
type: kubernetes.io/tls
```

Note: While it is possible to generate the values for `tls.crt` and `tls.key` manually, the above yaml can be generated using the following command:

```bash
kubectl create secret tls -n nlk nlk-tls-ca-secret --cert=/tmp/ca.crt --key=/tmp/ca.key --dry-run=client -o yaml > ca-secret.yaml
kubectl create secret tls -n nlk nlk-tls-server-secret --cert=/tmp/server.crt --key=/tmp/server.key --dry-run=client -o yaml > server-secret.yaml
kubectl create secret tls -n nlk nlk-tls-client-secret --cert=/tmp/client.crt --key=/tmp/client.key --dry-run=client -o yaml > client-secret.yaml
```

> [!WARNING]
> It is important that these files do not make their way into a public repository or other storage location where they can be accessed by unauthorized users.


Once the yaml files are generated they can be applied using `kubectl`:

```bash
kubectl apply -f ca-secret.yaml
kubectl apply -f server-secret.yaml
kubectl apply -f client-secret.yaml
```

# Verification

The Secrets can be verified using `kubectl`:

```bash
kubectl describe secret -n nlk nlk-tls-ca-secret
kubectl describe secret -n nlk nlk-tls-server-secret
kubectl describe secret -n nlk nlk-tls-client-secret
```

The output should look similar to the example above.

To see the actual values of the certificates, the following command can be used:

```bash
kubectl get secret -n nlk nlk-tls-ca-secret -o json | jq -r '.data["tls.crt"], .data["tls.key"]' | base64 -d
kubectl get secret -n nlk nlk-tls-server-secret -o json | jq -r '.data["tls.crt"], .data["tls.key"]' | base64 -d
kubectl get secret -n nlk nlk-tls-client-secret -o json | jq -r '.data["tls.crt"], .data["tls.key"]' | base64 -d
```

Note that this requires `jq` to be installed.

## References

- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [kubectl dry run flags](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#-em-dry-run-em-)