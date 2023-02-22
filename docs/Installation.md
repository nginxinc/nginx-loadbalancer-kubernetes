# Nginx Kubernetes Loadbalancer Solution

<br/>

## This is the Installation Guide for the Nginx Kubernetes Loadbalancer controller Solution.  It contains the detailed instructions for implementing the different components for the Solution.

<br/>

## Pre-Requisites

- Working kubernetes cluster, with admin privleges
- Running nginx-ingress controller, either OSS or Plus. This install guide follows the instructions for deploying an Nginx Ingress Controller here:  https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-manifests/
- Optional:  Demo application, this install guide uses the Nginx Cafe example, found here:  https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/ingress-resources/complete-example
- A bare metal Linux server or VM for the external LB Server, connected to a network external to the cluster.
- NginxPlus software loaded on the Edge Server, this install guide follows the instructions for installing NginxPlus on Centos 7, located here: https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-plus/
- The Nginx Kubernetes Loadbalancer (NKL) controller, new software for the Solution

<br/>

## Kubernetes Cluster

A standard K8s cluster is all that is required.  There must be enough resources to run the Nginx Ingress Controller, and the Nginx Kubernetes Loadbalancer Controller.  You must have administrative access to be able to create the namespace, services, and deployments for this Solution.  This Solution was tested on Kubernetes version 1.23.  Most recent versions => v1.21 should work just fine.

<br/>

## Nginx Ingress Controller

This is not part of the actual Solution, but it is the destination target for traffic (north-south) that is being sent to the cluster.  The installation of the actual Ingress Controller is outside the scope of the installation guide, but we include the links to the docs for your reference.  The NIC installation must follow the documents exactly as written, as this Solution refers to the `nginx-ingress` namespace and service objects.  Only the very last step is changed.

>Important!  The very last step in the NIC deployment with manifests, is to deploy the NodePort.yaml Service file.  `This file must be changed!  It is not the default nodeport file provided.`  Use the `nodeport-nkl.yaml` manifest file that is provided here.  The "ports name" in the Nodeport manifest `MUST` be in the correct format for this Solution to work correctly.  The port name is the mapping from NodePort to the LB Server's upstream blocks.  The port names are intentionally changed to avoid conflicts with other NodePort definitions.

Review the new NodePort Service defintion file:

```yaml
# NKL Nodeport Service file
# NodePort port name must be in the format of
# nkl-<upstream-block-name>
# Chris Akker, Jan 2023
#
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
  namespace: nginx-ingress
spec:
  type: NodePort 
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: nkl-nginx-lb-http   # Must be changed
  - port: 443
    targetPort: 443
    protocol: TCP
    name: nkl-nginx-lb-https   # Must be changed
  selector:
    app: nginx-ingress

```


```bash

kubectl apply -f nodeport-nkl.yaml

```

<br/>

## Demo Application

This is not part of the actual Solution, but it is useful to have a well-known application running in the cluster, as a useful target for test commands.  The example provided here is used by the Solution to demonstrate K8s application health check monitoring, to determine if the application is running in the cluster.  If you choose a different Application to test with, the health checks provided will NOT work, and will need to be modified to work correctly.

- Do not use the `cafe-ingress.yaml` file.  Rather, use the `cafe-virtualserver.yaml` file that is provided here.  It uses the Nginx CRDs to define a VirtualServer, and the related Routes and Redirects needed.  The redirects are required for the LB Server's health checks to work correctly!

```yaml
#Example virtual server with routes for Cafe Demo
#For NKL Solution, redirects required for LB Server health checks
#Chris Akker, Jan 2023
#
apiVersion: k8s.nginx.org/v1
kind: VirtualServer
metadata:
  name: cafe-vs
spec:
  host: cafe.example.com
  tls:
    secret: cafe-secret
    redirect:
      enable: true  #Redirect from http > https
      code: 301
  upstreams:
  - name: tea
    service: tea-svc
    port: 80
    lb-method: round_robin
    slow-start: 20s
    healthCheck:
      enable: true
      path: /tea
      interval: 20s
      jitter: 3s
      fails: 5
      passes: 2
      connect-timeout: 30s
      read-timeout: 20s
  - name: coffee
    service: coffee-svc
    port: 80
    lb-method: round_robin
    healthCheck:
      enable: true
      path: /coffee
      interval: 10s
      jitter: 3s
      fails: 3
      passes: 2
      connect-timeout: 30s
      read-timeout: 20s
  routes:
  - path: /
    action:
      redirect:
        url: https://cafe.example.com/coffee
        code: 302  #Redirect from / > /coffee
  - path: /tea
    action:
      pass: tea
  - path: /coffee
    action:
      pass: coffee
```

<br/>

## Linux VM or bare-metal LB Server

This is a standard Linux OS system, based on the Linux Distro and Technical Specs required for NginxPlus, which can be found here: https://docs.nginx.com/nginx/technical-specs/   This installation guide followed the "Installation of Nginx Plus on Centos/Redhat/Oracle" steps for installing Nginx Plus.  Note:  This solution will not work with Nginx OpenSource, as OpenSource does not have the API that is used in this Solution.  Installation on unsupported Distros is not Supported.

<br/>

## NginxPlus LB Server

This is the configuration required for the LB Server, external to the cluster.  It must be configured for the following.

- Move the Nginx default Welcome page from port 80 to port 8080.  Port 80 will be used by the stream context, instead of the http context.
- API write access enabled
- Plus Dashboard enabled, used for testing and visualization of the solution working
- Stream context enabled, for TCP loadbalancing
- Stream TCP loadbalancing configured

After a new installation of Nginx Plus, make the following configuration changes:

- Change Nginx's http default server to port 8080.  See the included `default.conf` file.  After reloading nginx, the default `Welcome to Nginx` page will be located at http://localhost:8080.

- Use the dashboard.conf file provided.  It will enable the /api endpoint, change the port to 9000, and provide access to the Plus dashboard.  Place this file in the /etc/nginx/conf.d folder, and reload nginx.  The Plus dashboard is now accessible at <server-ip>:9000/dashboard.html.  It should look similar to this:

< ss here>

- Create a new folder for the stream config files.  /etc/nginx/stream was used in this Solution.

```bash
mkdir /etc/nginx/stream
```

- Create 2 new `STATE` files for Nginx.  These are used to backup the configuration, in case Nginx restarts/reloads.

  Nginx State Files Required for Upstreams
    state file /var/lib/nginx/state/nginx-lb-http.state
    state file /var/lib/nginx/state/nginx-lb-https.state

```bash
touch /var/lib/nginx/state/nginx-lb-http.state
touch /var/lib/nginx/state/nginx-lp-https.state
```

- Enable the `stream` context for Nginx, which provides TCP load balancing.  See the included nginx.conf file.  Notice that the stream context is no longer commented out, a new folder is included, and a new stream.log logfile is used to track requests/responses.

- Configure Nginx Stream for TCP loadbalancing for this Solution.  Place this file in the /etc/nginx/stream folder.

```bash
# NginxK8sLB Stream configuration, for L4 load balancing
# Chris Akker, Jan 2023
# TCP Proxy and load balancing block
# Nginx Kubernetes Loadbalancer
# State File for persistent reloads/restarts
# Health Check Match example for cafe.example.com
#
#### nginxk8slb.conf

   upstream nginx-lb-http {
      zone nginx-lb-http 256k;
      state /var/lib/nginx/state/nginx-lb-http.state; 
    }

   upstream nginx-lb-https {
      zone nginx-lb-https 256k;
      state /var/lib/nginx/state/nginx-lb-https.state; 
    }

   server {
      listen 80;
      status_zone nginx-lb-http;
      proxy_pass nginx-lb-http;
      health_check match=cafe;
    }
             
   server {
      listen 443;
      status_zone nginx-lb-https;
      proxy_pass nginx-lb-https;
      health_check match=cafe;
    }

   match cafe {
      send "GET cafe.example.com/ HTTP/1.0\r\n";
      expect ~ "30*";
    }

```

<br/>

## Nginx Kubernetes Loadbalancing Controller

This is the new Controller, which is configured to watch the k8s environment, and send API updates to the Nginx LB Server when there are changes.  It only requires three things.
- New kubernetes namespace
- NKL Deployment, to start the Controller
- NKL ConfigMap, to configure the Controller


