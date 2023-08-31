# Kubernetes Secrets

## Overview

Kubernetes Secrets are used to provide the required certificates to NLK. There are two ways to create the Secrets:
- Using `kubectl`
- Using yaml files

The filenames for the certificates created previously are required for both methods. The examples below assume the certificates were generated in `/tmp` 
and follow the naming conventions in the documentation. 

## Using `kubectl`

The easiest way to create the secret is by using `kubectl`:

```bash
kubectl create secret tls nlk-tls-ca-secret --cert=/tmp/{ca,server,client}.crt --key=/tmp/{ca,server,client}.key
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
  name: nlk-tls-client-secret
type: kubernetes.io/tls
```

Note: While it is possible to generate the values for `tls.crt` and `tls.key` manually, the above yaml can be generated using the following command:

```bash
kubectl create secret tls nlk-tls-client-secret --cert=/tmp/{ca,server,client}.crt --key=/tmp/{ca,server,client}.key --dry-run=client -o yaml > {ca,server,client}-secret.yaml
```

Once the yaml files are generated they can be applied using `kubectl`:

```bash
kubectl apply -f {ca,server,client}-secret.yaml
```

**Warning: it is important that these files do not make their way into a public repository or other storage location where they can be accessed by unauthorized users.**

## References

- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [kubectl dry run flags](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#-em-dry-run-em-)