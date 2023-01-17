#!/bin/bash

echo "Unapplying all RBAC resources..."

kubectl delete -f ServiceAccount.yaml
kubectl delete -f ClusterRole.yaml
kubectl delete -f ClusterRoleBinding.yaml
kubectl delete -f Secret.yaml
