#!/bin/bash

echo "Applying all RBAC resources..."

kubectl apply -f ServiceAccount.yaml
kubectl apply -f ClusterRole.yaml
kubectl apply -f ClusterRoleBinding.yaml
kubectl apply -f Secret.yaml
