
## NLK and Kuberentes MultiCluster Load Balancing with HTTP/S

<br/>

![Kubernetes](../media/kubernetes-icon.png) | ![NLK](../media/nlk-logo.png)| ![NGINX Plus](../media/nginx-plus-icon.png)
--- | --- | ---

<br/>

## Overview

<br/>

With the NGINX Plus Servers located external to the Cluster, using NGINX's advanced HTTP/S features provide Enterprise class traffic management Solutions.  Using NGINX in HTTP mode allows for many solutions, some of are:
  
- MultiCluster Active/Active Load Balancing
- Horizontal Cluster Scaling
- HTTP Split Clients - for `A/B, Blue/Green, and Canary` test and production traffic steering.  Allows Cluster operations/maintainence like:
  - Node upgrades / additions
  - Software upgrades/security patches
  - Cluster resource expansions - memory, compute, storage, network, nodes
  - Troubleshooting, `using Live Traffic if needed`
  - ^^ With NO downtime or reloads
- NGINX Zone Sync of KeyVal data
- API Gateway functions
- Advanced TLS Processing - MutualTLS, OCSP, FIPS, dynamic cert loading
- Advanced Security features - App Protect WAF Firewall, Oauth, OIDC/JWT, Dynamic Rate and Bandwidth limits, GeoIP, IP block/allow lists
- NGINX Java Script (NJS) for custom solutions

<br/>

## Reference Diagram for NGINX NLK HTTP MultiCluster Load Balancing Solution

<br/>

Multiple K8s Clusters, HA NGINX Plus LB Servers, NLK Controllers

![NLK MultiCluster Diagram](../media/nlk-multicluster-config.png)


<br/>

NLK Controller watching nginx-ingress Service and Updating HTTP Upstreams; for Service Type Loadbalancer or NodePort:

![NLK MultiCluster LoadBalancer](../media/nlk-cluster1-add-loadbalancer.png)
or
![NLK MultiCluster NodePort](../media/nlk-cluster1-add-nodeport.png)

<br/>

MultiCluster Load Balancing

![NLK MultiCluster Dashboard](../media/nlk-multicluster-upstreams.png)

<br/>

NGINX HTTP Split Clients with Dynamic Ratio -- 10% Cluster1 : 90% Cluster2 

![NGINX HTTP Split 10](../media/nlk-clusters-10.png)

- NGINX Zone Sync of KeyVal data
- API Gateway functions
- Advanced TLS Processing - MutualTLS, OCSP, FIPS, dynamic cert loading
- Advanced Security features - App Protect WAF Firewall, Oauth, JWT, Dynamic Rate and Bandwidth limits, GeoIP, IP block/allow lists
- NGINX Java Script (NJS) for custom solutions

<br/>

## Reference Diagram for NLK HTTP MultiCluster Load Balancing Solution

<br/>

Multiple K8s Clusters, HA NGINX Plus LB Servers, NLK Controllers

![NLK MultiCluster Diagram](../media/nlk-multicluster-config.png)


<br/>

NLK Watching nginx-ingress Service and Updating HTTP Upstreams; Service Type Loadbalancer or NodePort:

![NLK MultiCluster LoadBalancer](../media/nlk-cluster1-add-loadbalancer.png)
or
![NLK MultiCluster NodePort](../media/nlk-cluster1-add-nodeport.png)

<br/>

MultiCluster Load Balancing

![NLK MultiCluster Dashboard](../media/nlk-multicluster-upstreams.png)

<br/>

NGINX HTTP Split Clients with Dynamic Ratio -- 10% Cluster1 : 90% Cluster2 

![NGINX Grafana Split 10](../media/nlk-clusters-10.png)

<br/>

### Grafana Charts - Examples showing 10, 90, 50% Split Ratios

<br/>

NGINX HTTP Split Clients with Dynamic Ratio -- 10% Cluster1 : 90% Cluster2 

![NGINX HTTP Split 10](../media/nlk-grafana-reqs-10.png)

<br/>

NGINX HTTP Split Clients with Dynamic Ratio -- 90% Cluster1 : 10% Cluster2 

![NGINX Grafana Split 10](../media/nlk-grafana-reqs-90.png)
<br/>

NGINX HTTP Split Clients with Dynamic Ratio -- 50% Cluster1 : 50% Cluster2 

![NGINX Grafana Split 10](../media/nlk-grafana-reqs-50.png)

<br/>

The `Installation Guide` for HTTP MultiCluster Solution is located in the docs/http folder:

[HTTP MultiCluster Loadbalancing Guide](../http/http-installation-guide.md)

<br/>

## Authors
- Chris Akker - Solutions Architect - Community and Alliances @ F5, Inc.
- Steve Wagner - Solutions Architect - Community and Alliances @ F5, Inc.
