#!/bin/bash

echo "Unapplying all RBAC resources..."

kubectl delete -f serviceaccount.yaml
kubectl delete -f clusterrole.yaml
kubectl delete -f clusterrolebinding.yaml
kubectl delete -f secret.yaml
