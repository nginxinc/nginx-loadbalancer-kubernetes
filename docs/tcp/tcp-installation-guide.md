# NGINX Loadbalancer for Kubernetes Solution

<br/>

## This is the `TCP Installation Guide` for the NGINX Loadbalancer for Kubernetes Controller Solution.  It contains detailed instructions for implementing the different components for the Solution.

<br/>

![Kubernetes](../media/kubernetes-icon.png) | ![NLK](../media/nlk-logo.png) | ![NGINX Plus](../media/nginx-plus-icon.png) | ![NIC](../media/nginx-ingress-icon.png)
--- | --- | --- | ---

<br/>

## Solution Overview

### This Solution from NGINX provides Enterprise class features which address common challenges with networking, traffic management, and High Availability for On-Premises Kubernetes Clusters.

1. Provides a `replacement Loadbalancer Service.`  The Loadbalancer Service is a key component provided by most Cloud Providers.  However, when running a cluster On Premises, the `Loadbalancer Service is not available`.  This Solution provides a replacement, using an NGINX Server, and a new K8s Controller.  These two components work together to watch the NodePort Service in the cluster, and immediately update the NGINX Loadbalancing Server when changes occur.  No more static `ExternalIP` needed in your `loadbalancer.yaml` Manifests!
2. Provides automatic NGINX upstream config updates, application health checks, advanced Loadbalancing algorithms, and enhanced metrics.
3. Provides an upgrade option to NGINX's powerful HTTP processing - `dynamic, ratio-based Load Balancing for Multiple Clusters.`  This allows for advanced traffic steering, and operation efficiency with no Reloads or downtime.  See the HTTP Install Guide for additional details on the advanced HTTP Solution, which can provide:
  - MultiCluster Loadbalancing and High Availability
  - Horizontal Cluster scaling
  - Non-stop seemless K8s Cluster upgrades, migrations, patching
  - HTTP Split clients for `A/B, Blue/Green, and Canary testing` and production traffic
  - Additional security features like App Protect Firewall, JWT auth, Rate Limiting, Service and Bandwidth controls, FIPS, advanced TLS features.

<br/>

![NLK Solution Overview](../media/nlk-stream-diagram.png)

<br/>

## Installation Steps

1. Install NGINX Ingress Controller in your Cluster
2. Install NGINX Cafe Demo Application in your Cluster
3. Install NGINX Plus on the Loadbalancer Server(s)
4. Configure NGINX Plus for TCP Load Balancing 
5. Install NLK NGINX Loadbalancing for Kubernetes Controller in your Cluster
6. Install NLK LoadBalancer or NodePort Service manifest
7. Test out NLK

<br/>

### Pre-Requisites

1. Working Kubernetes cluster, with admin privleges
   
2. Running `nginx-ingress controller`, either OSS or Plus. This install guide followed the instructions for deploying an NGINX Ingress Controller here:  https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-manifests
   
3. Demo application, this install guide uses the NGINX Cafe example, found here:  https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/ingress-resources/complete-example
   
4. A bare metal Linux server or VM for the external NGINX Loadbalancing Server, connected to a network external to the cluster.  Two of these will be required if High Availability is needed, as shown here.
   
5. NGINX Plus software loaded on the Loadbalancing Server(s). This install guide follows the instructions for installing NGINX Plus on Centos 7, located here: https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-plus/
   
6. The NGINX Loadbalancer for Kubernetes (NLK) Controller, new software from NGINX for this Solution.

<br/>

## Kubernetes Cluster

<br/>

![Kubernetes](../media/kubernetes-icon.png)

<br/>

A standard K8s cluster is all that is required.  There must be enough resources available to run the NGINX Ingress Controller, and the new NGINX Loadbalancer for Kubernetes Controller, and test application like the Cafe Demo.  You must have administrative access to be able to create the namespace, services, and deployments for this Solution.  This Solution was tested on Kubernetes version 1.23.

<br/>

## 1. Install NGINX Ingress Controller

<br/>

![NIC](../media/nginx-ingress-icon.png)

<br/>

The NGINX Ingress Controller in this Solution is the destination target for traffic (north-south) that is being sent to the cluster.  The installation of the actual Ingress Controller is outside the scope of this installation guide, but we include the links to the docs for your reference.  The NIC installation must follow the documents exactly as written, as this Solution refers to the `nginx-ingress` namespace and service objects.  **Only the very last step is changed.**  

**NOTE:** This Solution only works `with nginx-ingress` from NGINX.  It will not work with the Community version of Ingress, called ingress-nginx.  

If you are unsure which Ingress Controller you are running, check out the blog on nginx.com:  
https://www.nginx.com/blog/guide-to-choosing-ingress-controller-part-4-nginx-ingress-controller-options

<br/>

>Important!  Do not complete the very last step in the NIC deployment with Manifests, `do not deploy the loadbalancer.yaml or nodeport.yaml Service file!`  You will apply a different loadbalancer or nodeport Service manifest later, after the NLK Controller is up and running.  `The nginx-ingress Service file must be changed` - it is not the default file. 

<br/>

## 2. Install NGINX Ingress Demo Application

<br/>

![Cafe Dashboard](../media/cafe-dashboard.png)

<br/>

This is not part of the actual Solution, but it is useful to have a well-known application running in the cluster, as a known-good target for test commands.  The example provided here is used by the Solution to demonstrate proper traffic flows.  

Note: If you choose a different Application to test with, `the NGINX health checks provided here will likely NOT work,` and will need to be modified to work correctly.

<br/>

1. Use the provided Cafe Demo manifests in the cafe-demo folder:

    ```bash
    kubectl apply -f cafe-secret.yaml
    kubectl apply -f cafe.yaml
    kubectl apply -f cafe-virtualserver.yaml
    ```

1. The Cafe Demo reference files are located here:

    https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/ingress-resources/complete-example

1. The Cafe Demo Docker image used here is an upgraded one, with simple graphics and additional TCP/IP and HTTP variables added.

    https://hub.docker.com/r/nginxinc/ingress-demo

    **IMPORTANT** - Do not use the `cafe-ingress.yaml` file.  Rather, use the `cafe-virtualserver.yaml` file that is provided here.  It uses the NGINX Plus CRDs to define a VirtualServer, and the related Virtual Server Routes needed.  If you are using NGINX OSS Ingress Controller, you will need to use the appropriate manifests, which is not covered in this Solution.

<br/>

## 3. Install NGINX Plus on LoadBalancing Server(s)

<br/>

### Linux VM or bare-metal Loadbalancing Server

![Linux](../media/linux-icon.png) | ![NGINX Plus](../media/nginx-plus-icon.png)
--- | ---

<br/>

This can be any standard Linux OS system, based on the Linux Distro and Technical Specs required for NGINX Plus, which can be found here: https://docs.nginx.com/nginx/technical-specs/   

This Solution followed the `Installation of NGINX Plus on Centos/Redhat/Oracle` steps for installing NGINX Plus.  

https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-plus/

>NOTE:  This Solution will only work with NGINX Plus, as NGINX OpenSource does not have the API that is used in this Solution.  Installation on unsupported Linux Distros is not recommended.

If you need a license for NGINX Plus, a 30-day Trial license is available here:

https://www.nginx.com/free-trial-request/

<br/>

## 4. Configure NGINX Plus for TCP Load Balancing

<br/>

### This is the configuration required for the NGINX Loadbalancing Server, external to the cluster.  It must be configured for the following:

1. Move the NGINX default Welcome page from port 80 to port 8080.  Port 80 will be used by the stream context, instead of the http context.
   
2. Plus API with write access enabled on port 9000.
   
3. Plus Dashboard enabled, used for testing, monitoring, and visualization of the Solution working.
   
4. The NGINX `stream` context is enabled, and configured for TCP loadbalancing.

<br/>

### Overview of the Config Files used for the NGINX Plus Loadbalancing Servers:

<br/>

For easy installation/configuration, Git Clone this repository onto the Loadbalancing Server, it contains all the example files that are used here.

```bash
https://github.com/nginxinc/nginx-loadbalancer-kubernetes.git
```

<br/>


```bash
etc/
└── nginx/
    ├── conf.d/
    │   ├── dashboard.conf........ NGINX Plus API and Dashboard config
    │   └── default.conf.......... New default.conf config
    ├── nginx.conf................ New nginx.conf
    └── stream
        └── nginxk8slb.conf....... NGINX TCP Loadbalancing config 
```

``` bash
nginx-loadbalancer-kubernetes/
└── docs/
    └── tcp/
        ├── loadbalancer-nlk.yaml........ LoadBalancer manifest
        └── nodeport-nlk.yaml ........... NodePort manifest
```
<br/>

After the new installation of NGINX Plus, make the following configuration changes:

1. Change NGINX's http default server to port 8080.  See the included `default-tcp.conf` file.  After reloading NGINX, the default `Welcome to NGINX` page will be located at http://localhost:8080.

    ```bash
    cat /etc/nginx/conf.d/default.conf
    # NGINX Loadbalancing for Kubernetes Solution
    # Chris Akker, Apr 2023
    # Example default.conf
    # Change default_server to port 8080
    #
    server {
        listen       8080 default_server;   # Changed to 8080
        server_name  localhost;

        #access_log  /var/log/nginx/host.access.log  main;

        location / {
            root   /usr/share/nginx/html;
            index  index.html index.htm;
        }

        #error_page  404              /404.html;

        # redirect server error pages to the static page /50x.html
        #
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   /usr/share/nginx/html;
        }

    ### other sections removed for clarity

    }

    ```

2. Enable the NGINX Plus dashboard.  Use the `dashboard.conf` file provided.  It will enable the /api endpoint, change the port to 9000, and provide access to the Plus Dashboard.  Note:  There is no security for the /api endpoint in this example config, it should be secured as approprite with TLS or IP allow list.

    Place this file in the /etc/nginx/conf.d folder, and reload nginx.  The Plus dashboard is now accessible at http://nginx-lbserver-ip:9000/dashboard.html.  It should look similar to this:

    ![NGINX Dashboard](../media/nlk-stream-dashboard.png)

3. Create a new folder for the NGINX stream .conf files.  `/etc/nginx/stream` is used in this Solution.

    ```bash
    mkdir /etc/nginx/stream
    ```

4. Enable the `stream` context for NGINX, which provides TCP load balancing.  See the included nginx.conf file.  Notice that the stream context is no longer commented out, the new folder is included, and a new stream.log logfile is used to track requests/responses.

    ```bash
    cat /etc/nginx/nginx.conf

    # NGINX Loadbalancing for Kubernetes Solution
    # Chris Akker, Apr 2023
    # Example nginx.conf
    # Enable Stream context, add /var/log/nginx/stream.log
    #

    user  nginx;
    worker_processes  auto;

    error_log  /var/log/nginx/error.log notice;
    pid        /var/run/nginx.pid;

    events {
        worker_connections  1024;
    }

    http {
        include       /etc/nginx/mime.types;
        default_type  application/octet-stream;

        log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                          '$status $body_bytes_sent "$http_referer" '
                          '"$http_user_agent" "$http_x_forwarded_for"';

        access_log  /var/log/nginx/access.log  main;

        sendfile        on;
        #tcp_nopush     on;

        keepalive_timeout  65;

        #gzip  on;

        include /etc/nginx/conf.d/*.conf;
    }

    # TCP/UDP proxy and load balancing block
    #
    stream {

      include  /etc/nginx/stream/*.conf;

      log_format  stream  '$remote_addr - $server_addr [$time_local] $status $upstream_addr $upstream_bytes_sent';

      access_log  /var/log/nginx/stream.log  stream;
    }

    ```

5. Configure NGINX Stream for TCP loadbalancing for this Solution.

    `Notice that this example Solution uses Ports 80 and 443.`  
  
    Place this file in the /etc/nginx/stream folder, and reload NGINX.  Notice the match block and health check directives are for the cafe.example.com Demo application from NGINX.

    ```bash
    # NGINX Loadbalancing for Kubernetes Stream configuration, for L4 load balancing
    # Chris Akker, Apr 2023
    # TCP Proxy and load balancing block
    # NGINX Loadbalancer for Kubernetes
    # State File for persistent reloads/restarts
    # Health Check Match example for cafe.example.com
    #
    #### nginxk8slb.conf

      upstream nginx-lb-http {
          zone nginx-lb-http 256k;
          #servers managed by NLK Controller
          state /var/lib/nginx/state/nginx-lb-http.state; 
        }

      upstream nginx-lb-https {
          zone nginx-lb-https 256k;
          #servers managed by NLK Controller
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

6. Check the NGINX Plus Dashboard, at http://nginx-lbserver-ip:9000/dashboard.html.  You should see something like this:

    ![NLK Stream Upstreams](../media/nlk-stream-dashboard.png)

7. If you have 2 NGINX Loadbalancing Servers for High Availability, repeat the previous NGINX Plus installation and configuration steps on the second Loadbalancing Server.

<br/>

## 5. Install NLK - NGINX Loadbalancing for Kubernetes Controller

<br/>

![NIC](../media/nlk-logo.png)

<br/>

### This is the new K8s Controller from NGINX, which is configured to watch the k8s environment, the `nginx-ingress` Service object, and send API updates to the NGINX Loadbalancing Server(s) when there are changes.  It only requires three things:

1. New Kubernetes namespace and RBAC
2. NLK ConfigMap, to configure the Controller
3. NLK Deployment, to deploy and run the Controller

<br/>

- Create the new K8s namespace:

```bash
kubectl create namespace nlk
```

- Apply the manifests for Secret, Service, ClusterRole, and ClusterRoleBinding:

```bash
kubectl apply -f secret.yaml serviceaccount.yaml clusterrole.yaml clusterrolebinding.yaml
```

Modify the ConfigMap manifest to match your NGINX Loadbalancing Server(s). Change the `nginx-hosts` IP address to match your NGINX Loadbalancing Server IP.  If you have 2 or more Loadbalancing Servers, separate them with a comma.  Keep the port number for the Plus API endpoint, and the `/api` URL as shown.

```yaml

apiVersion: v1
kind: ConfigMap
data:
  nginx-hosts:
    "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api"    # change IP(s) to match NGINX Loadbalancing Server(s)
metadata:
  name: nlk-config
  namespace: nlk

```

Apply the updated ConfigMap:

```bash
kubectl apply -f nlk-configmap.yaml
```

Deploy the NLK Controller:

```bash
kubectl apply -f nlk-deployment.yaml
```

Check to see if the NLK Controller is running, with the updated ConfigMap:

```bash
kubectl get pods -n nlk
```
```bash
kubectl describe cm nlk-config -n nlk
```

The status should show "running", your `nginx-hosts` should have the Loadbalancing Server IP:9000/api.

![NLK Running](../media/nlk-configmap.png)

To make it easy to watch the NLK Controller's log messages, add the following bash alias:

```bash
alias nlk-follow-logs='kubectl -n nlk get pods | grep nlk-deployment | cut -f1 -d" "  | xargs kubectl logs -n nlk --follow $1'
```

Using a Terminal, you can watch the NLK Controller log:

```bash
nlk-follow-logs
```

Leave this Terminal window open, so you can watch the log messages.

<br/>

## 6. Install NLK Loadbalancer or NodePort Service Manifest

<br/>

Select which Service Type you would like, and follow the appropriate steps below.  Do not use both the LoadBalancer and NodePort Service files at the same time.

Instead, use the `loadbalancer-nlk.yaml` or `nodeport-nlk.yaml` manifest file that is provided here with this Solution.  The "ports name" in the manifests `MUST` be in the correct format for this Solution to work correctly.  
>**`The port name is the mapping from NodePorts to the Loadbalancing Server's upstream blocks.`**  The port names are intentionally changed to avoid conflicts with other NodePort definitions.

<br/>

### If you want to run a Service Type LoadBalancer

Review the new `loadbalancer-nlk.yaml` Service definition file:

```yaml
# NLK LoadBalancer Service file
# Spec -ports name must be in the format of
# nlk-<upstream-block-name>
# The nginxinc.io Annotation must be added
# externalIPs are set to NGINX Loadbalancing Servers
# Chris Akker, Apr 2023
#
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
  namespace: nginx-ingress
  annotations:
    nginxinc.io/nlk-nginx-lb-http: "stream"    # Must be added
    nginxinc.io/nlk-nginx-lb-https: "stream"   # Must be added
spec:
  type: LoadBalancer
  externalIPs:
  - 10.1.1.4          #NGINX Loadbalancing1 Server
  - 10.1.1.5          #NGINX Loadbalancing2 Server
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: nlk-nginx-lb-http      # Must be changed
  - port: 443
    targetPort: 443
    protocol: TCP
    name: nlk-nginx-lb-https     # Must be changed
  selector:
    app: nginx-ingress

```

- Apply the NLK Compatible LoadBalancer `loadbalancer-nlk.yaml` Service Manifest:

```bash
kubectl apply -f loadbalancer-nlk.yaml
```

- Verify the LoadBalancer is now defined:

```bash
kubectl get svc nginx-ingress -n nginx-ingress
```

The nginx-ingress Service, `ExternalIPs` should match your external NGINX Loadbalancing Server IP(s):

![NLK Stream Loadbalancer](..//media/nlk-stream-add-loadbalancer.png)

Legend:
- Orange is the `TYPE LoadBalancer` Service.
- Red is the LoadBalancer Service `EXTERNAL-IP`, which are your NGINX Loadbalancing Server IP(s); 10.1.1.4 and 10.1.1.5 in this example.
- Blue is the `K8s NodePort mapping` for Port 80.
- Indigo is the `K8s NodePort mapping` for Port 443.
- Green is the NLK Log messages, creating the upstreams to match.
- The new NLK Controller updates the NGINX Loadbalancing Server upstreams with these, shown on the dashboard.

>>No Reload of NGINX needed!  The NLK Controller uses the Plus API to dynamically add/delete/modify the upstreams as the `nginx-ingress Service` changes.

<br/>

### Alternatively, if you want a Service Type NodePort

Review the new `nodeport-nlk.yaml` Service defintion file:

```yaml
# NLK Nodeport Service file
# NodePort -ports name must be in the format of
# nlk-<upstream-block-name>
# The nginxinc.io Annotation must be added
# Chris Akker, Apr 2023
#
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
  namespace: nginx-ingress
  annotations:
    nginxinc.io/nlk-nginx-lb-http: "stream"    # Must be added
    nginxinc.io/nlk-nginx-lb-https: "stream"   # Must be added
spec:
  type: NodePort 
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
    name: nlk-nginx-lb-http     # Must be changed
  - port: 443
    targetPort: 443
    protocol: TCP
    name: nlk-nginx-lb-https    # Must be changed
  selector:
    app: nginx-ingress

```

- Create the NLK compatible NodePort Service, using the `nodeport-nlk.yaml` manifest provided:

```bash
kubectl apply -f nodeport-nlk.yaml
```

- Verify the NodePort is now defined:

```bash
kubectl get svc nginx-ingress -n nginx-ingress
```

![NLK NodePort](../media/nlk-stream-nodeport.png)
![NLK Stream Upstreams Dashboard](../media/nlk-stream-upstreams.png)

Legend:
- Orange is the `TYPE NodePort` Service.
- Notice the EXTERNAL-IP is blank, as expected.
- Blue is the `K8s NodePort mapping` for Port 80.
- Indigo is the `K8s NodePort mapping` for Port 443.

### NodePort mapping is 80:31681 and 443:31721,  K8s Workers are 10.1.1.8 and .10.

<br/>

### Deep Dive Explanation

<br/>

The name of the Service port is matched to the name of the upstream block in NGINX.  The Plus API, follows a defined format, so the url for the API call must be correct, in order to update the correct NGINX upstream block.  There are 2 types of upstreams in NGINX.  `Stream` upstreams are used in the stream context, for TCP/UDP load balancing configurations.  `Http` upstreams are used in the http context, for HTTP/HTTPS configurations.  (See details for HTTP in the http-installation-guide.md, here:  [HTTP Guide](../http/http-installation-guide.md).

<br/>

## 7. Testing NLK NGINX Loadbalancer for Kubernetes

<br/>

When you are finished, the NGINX Plus Dashboard on the Loadbalancing Server should look similar to the following image:

![NGINX Upstreams Dashboard](../media/nlk-stream-upstreams.png)

Important items for reference:
- Orange are the upstream server blocks, from the `etc/nginx/stream/nginxk8slb.conf` file.
- Blue is the IP:Port of the nginx-ingress Service for http.
- Indigo is the IP:Port of the nginx-ingress Service for https.

>Note: In this example, there is a 3-Node K8s cluster, with one Control Node, and 2 Worker Nodes.  The NLK Controller only configures `Worker Node` IP addresses, which are:
- 10.1.1.8
- 10.1.1.10

Note:  K8s Control Nodes are excluded intentionally.

<br/>

Configure DNS, or your local hosts file, for cafe.example.com > NGINXLoadbalancing Server IP Address.  In this example:


```bash
cat /etc/hosts
10.1.1.4 cafe.example.com
```

- Open a browser tab to https://cafe.example.com/coffee.

The Dashboard's `TCP/UDP Upstreams Connection counters` will increase as you refresh the browser page several times.

- Using a Terminal, delete the `nginx-ingress loadbalancer service` or `nginx-ingress nodeport service` definition. 

```bash
kubectl delete -f loadbalancer-nlk.yaml
```
or
```bash
kubectl delete -f nodeport-nlk.yaml
```

Now the `nginx-ingress` Service is gone, and the upstream lists will be empty in the Dashboard.

![NGINX No NodePort](../media/nlk-stream-no-nodeport.png)

The NLK log messages confirm the deletion of the upstreams:

![NLK Logs Deleted](../media/nlk-stream-logs-deleted.png)

- If you refresh the cafe.example.com browser page, it will Time Out.  There are NO upstreams for NGINX to send the request to!

---

- Add the `nginx-ingress` Service back to the cluster:

```bash
kubectl apply -f loadbalancer-nlk.yaml
```
or
```bash
kubectl apply -f nodeport-nlk.yaml
```

- Verify the nginx-ingress Service is re-created.  Notice the the NodePort Numbers have changed!

```bash
kubectl get svc nginx-ingress -n nginx-ingress
```

`The NLK Controller detects this change, and modifies the Loadbalancing Server(s) upstreams to match.`  The Dashboard will show you the new Port numbers, matching the new LoadBalancer or NodePort definitions.  The NLK logs show these messages, confirming the changes:

![NLK LoadBalancer](../media/nlk-stream-add-loadbalancer.png)

or

![NLK NodePort](../media/nlk-stream-nodeport.png)
![NLK Logs Created](../media/nlk-stream-logs-created.png)
![NGINX Upstreams Dashboard](../media/nlk-stream-upstreams.png)

<br/>

The Completes the Testing Section.

<br/>

## Authors
- Chris Akker - Solutions Architect - Community and Alliances @ F5, Inc.
- Steve Wagner - Solutions Architect - Community and Alliances @ F5, Inc.

