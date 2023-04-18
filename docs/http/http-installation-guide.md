# NGINX Kubernetes Loadbalancer - MultiCluster LB Solution

<br/>

## This is the `HTTP Installation Guide` for the NGINX Kubernetes Loadbalancer Controller Solution.  It contains detailed instructions for implementing the different components for the Solution.

<br/>

![Kubernetes](../media/kubernetes-icon.png) | ![NKL](../media/nkl-logo.png)| ![NGINX Plus](../media/nginx-plus-icon.png) | ![NIC](../media/nginx-ingress-icon.png)
--- | --- | --- | ---

<br/>

## Solution Overview

### This Solution from NGINX provides Enterprise class features which address common challenges with networking, traffic management, and High Availability for On-Premises Kubernetes Clusters.

1. Provides a `replacement Loadbalancer Service.`  The Loadbalancer Service is a key component provided by most Cloud Providers.  However, when running a K8s Cluster On Premises, the `Loadbalancer Service is not available.`  
This Solution provides a replacement, using an NGINX Server, and a new K8s Controller from NGINX.  These two components work together to watch the `NodePort Service` in the cluster, and immediately update the NGINX LB Server when changes occur.  
>No static `ExternalIP` needed in your loadbalancer.yaml Manifests!
2. Provides `MultiCluster Load Balancing`, traffic steering, health checks, TLS termination, advanced LB algorithms, and enhanced metrics.
3. Provides dynamic, ratio-based Load Balancing for Multiple Clusters.  This allows for advanced traffic steering, and operation efficiency with no Reloads or downtime.
   - MultiCluster Active/Active Load Balancing
   - Horizontal Cluster Scaling
   - HTTP Split Clients - for A/B, Blue/Green, and Canary test and production traffic steering.  Allows Cluster operations/maintainence like upgrades, patching, expansion and troubleshooting
   - NGINX Zone Sync of KeyVal data
   - Advanced TLS Processing - MutualTLS, OCSP, FIPS, dynamic cert loading
   - Advanced Security features - App Protect WAF Firewall, Oauth, JWT, Dynamic Rate and Bandwidth limits, GeoIP, IP block/allow lists
   - NGINX Java Script (NJS) for custom solutions

<br/>

![NKL Solution Overview](../media/nkl-multicluster-diagram.png)

<br/>

## Installation Steps

1. Install NGINX Ingress Controller in your Cluster
2. Install NGINX Cafe Demo Application in your Cluster
3. Install NGINX Plus on the Loadbalancer Server(s) 
4. Configure NGINX Plus for MultiCluster Load Balancing
5. Install NKL - NGINX Kubernetes LB Controller in your Cluster
6. Test out NKL
7. Test MultiCluster Load Balancing Solution
8. Optional - Monitor traffic with Prometheus / Grafana

<br/>

### Pre-Requisites

- Working kubernetes clusters, with admin privleges
- Running `nginx-ingress controller`, either OSS or Plus. This install guide followed the instructions for deploying an NGINX Ingress Controller here:  https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-manifests/
- Demo application, this install guide uses the NGINX Cafe example, found here:  https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/ingress-resources/complete-example
- A bare metal Linux server or VM for the external NGINX LB Server, connected to a network external to the cluster.  Two of these will be required if High Availability is needed, as shown here.
- NGINX Plus software loaded on the LB Server(s). This install guide follows the instructions for installing NGINX Plus on Centos 7, located here: https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-plus/
- The NGINX Kubernetes Loadbalancer (NKL) Controller, new software for this Solution.

<br/>

### Kubernetes Clusters

<br/>

![Kubernetes](../media/kubernetes-icon.png)

<br/>

A standard K8s cluster is all that is required, two or more Clusters if you want the `Active/Active MultiCluster Load Balancing Solution` using HTTP Split Clients.  There must be enough resources available to run the NGINX Ingress Controller, and the NGINX Kubernetes Loadbalancer Controller, and test application like the Cafe Demo.  You must have administrative access to be able to create the namespace, services, and deployments for this Solution.  This Solution was tested on Kubernetes version 1.23.  Most recent versions => v1.21 should work just fine.

<br/>

## 1. Install NGINX Ingress Controller

<br/>

![NIC](../media/nginx-ingress-icon.png)

<br/>

The NGINX Ingress Controller in this Solution is the destination target for traffic (north-south) that is being sent to the cluster(s).  The installation of the actual Ingress Controller is outside the scope of this guide, but the links to the docs are included for your reference.  `The NIC installation using Manifests must follow the documents exactly as written,` as this Solution depends on the `nginx-ingress` namespace and service objects.  **Only the very last step is changed.**  

**NOTE:** This Solution only works with `nginx-ingress from NGINX`.  It will `not` work with the K8s Community version of Ingress, called ingress-nginx.  

If you are unsure which Ingress Controller you are running, check out the blog on nginx.com:  
    
https://www.nginx.com/blog/guide-to-choosing-ingress-controller-part-4-nginx-ingress-controller-options


>Important!  The very last step in the NIC deployment with Manifests, is to deploy the `nodeport.yaml` Service file.  `This file must be changed - it is not the default nodeport file.`  
Instead, use the `nodeport-cluster1.yaml` manifest file that is provided here with this Solution.  The "ports name" in the Nodeport manifest `MUST` be in the correct format for this Solution to work correctly.  The port name is the mapping from NodePorts to the LB Server's upstream blocks.  The port names are intentionally changed to avoid conflicts with other NodePort definitions.

Review the new `nodeport-cluster1.yaml` Service defintion file:

```yaml
# NKL Nodeport Service file
# Chris Akker, Apr 2023
# NodePort -ports name must be in the format of
#
## nkl-<upstream-block-name> ##
# 
#
apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
  namespace: nginx-ingress
spec:
  type: NodePort 
  ports:
  - port: 443
    targetPort: 443
    protocol: TCP
    name: nkl-cluster1-https  # This must match NGINX upstream name
  selector:
    app: nginx-ingress

```

Apply the updated nodeport-cluster1.yaml Manifest:

```bash
kubectl apply -f nodeport-cluster1.yaml
```

**NOTE:** If you have a second K8s cluster, and you want to Load Balance both Clusters using the MultiCluster Solution, repeat the appropriates steps on your second cluster.  

**IMPORTANT:  Do not mix and match nodeport-clusterX.yaml files.**  

- `nodeport-cluster1.yaml` must be used for Cluster1, `nodeport-cluster2.yaml` must be used for Cluster2.  The NodePort definitions must match each cluster exactly.
- Nodeports and manifest files must match the target cluster for the HTTP Split Clients dynamic ratio configuration to work correctly.  
- It is highly recommended that you configure, test, and verify traffic is flowing correctly on Cluster1, before you add Cluster2.  
- Be aware of and properly set your `./kube/config Cluster Context`, before applying one of these nodeport definitions.

<br/>

## 2. Install NGINX Cafe Demo Application

<br/>

![Cafe Dashboard](..//media/cafe-dashboard.png)

<br/>

This is not part of the actual Solution, but it is useful to have a well-known application running in the cluster, as a known-good target for test commands.  The example provided here is used by the Solution to demonstrate proper traffic flows.  

Note: If you choose a different Application to test with, `the NGINX health checks provided here will likely NOT work,` and will need to be modified to work correctly.

- Use the provided Cafe Demo manifests in the docs/cafe-demo folder:

  ```bash
  kubectl apply -f cafe-secret.yaml
  kubectl apply -f cafe.yaml
  kubectl apply -f cafe-virtualserver.yaml
  ```

- The Cafe Demo reference files are located here:

  https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/ingress-resources/complete-example

- The Cafe Demo Docker image used here is an upgraded one, with simple graphics and additional Request and Response variables added.

  https://hub.docker.com/r/nginxinc/ingress-demo

**IMPORTANT** - Do not use the `cafe-ingress.yaml` file.  Rather, use the `cafe-virtualserver.yaml` file that is provided here.  It uses the NGINX Plus CRDs to define a VirtualServer, and the related Virtual Server Routes needed.  If you are using NGINX OSS Ingress Controller, you will need to use the appropriate manifests, which is not covered in this Solution.

<br/>

## 3. Install NGINX Plus on LoadBalancer Server(s)

<br/>

### Linux VM or bare-metal LB Server

![Linux](../media/linux-icon.png) | ![NGINX Plus](../media/nginx-plus-icon.png)
--- | ---

<br/>

This is any standard Linux OS system, based on the Linux Distro and Technical Specs required for NGINX Plus, which can be found here: https://docs.nginx.com/nginx/technical-specs/   

This Solution followed the `Installation of NGINX Plus on Centos/Redhat/Oracle` steps for installing NGINX Plus.  

https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-plus/

>NOTE:  This solution will only work with NGINX Plus, as NGINX OpenSource does not have the API that is used in this Solution.  Installation on unsupported Linux Distros is not recommended.

If you need a license for NGINX Plus, a 30-day Trial license is available here:

https://www.nginx.com/free-trial-request/

<br/>

## 4. Configure Nginx Plus for MultiCluster Load Balancing

<br/>


### This is the NGINX configuration required for the NGINX LB Server, external to the cluster.  It must be configured for the following:

- Move the NGINX default Welcome page from port 80 to port 8080.  Port 80 will be used by Prometheus in this Solution.
- The NGINX NJS module is enabled, and configured to export the NGINX Plus statistics.
- A self-signed TLS cert/key are used in this example for terminating TLS traffic for the Demo application, https://cafe.example.com.
- Plus API with write access enabled on port 9000.
- Plus Dashboard enabled, used for testing, monitoring, and visualization of the Solution working.
- The `http` context is used for MultiCluster Loadbalancing, for HTTP/S processing, Split Clients ratio, and prometheus exporting.
- Plus KeyValue store is configured, to hold the dynamic Split ratio metadata.
- Plus Zone Sync on Port 9001 is configured, to synchronize the dynamic KVstore data between multiple NGINX LB Servers.

<br/>

- Overview of the Config Files used for the NGINX Plus LB Servers:

>/etc/nginx/conf.d

    - clusters.conf           | MultiCluster LB and split clients config
    - dashboard.conf          | NGINX Plus API and Dashboard config
    - default-http.conf       | New default.conf config
    - grafana-dashboard.json  | NGINX Plus Grafana dashboard
    - nginx.conf              | New nginx.conf
    - nodeport-cluster1.yaml  | NodePort config for Cluster1
    - nodeport-cluster2.yaml  | NodePort config for Cluster2
    - prometheus.conf         | NGINX Prometheus config
    - prometheus.yml          | Prometheus container config

>/etc/nginx/stream
       
    - zonesync.conf           | NGINX Zone Sync config

<br/>

After a new installation of NGINX Plus, make the following configuration changes:

- Change NGINX's http default server to port 8080.  See the included `default-http.conf` file.  After reloading nginx, the default `Welcome to NGINX` page will be located at http://localhost:8080.

```bash
cat /etc/nginx/conf.d/default.conf
# NGINX K8s Loadbalancer Solution
# Chris Akker, Jan 2023
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

- Enable the NGINX Plus dashboard.  Use the `dashboard.conf` file provided.  It will enable the /api endpoint, change the port to 9000, and provide access to the Plus Dashboard.  Note:  There is no security for the /api endpoint in this example config, it should be secured as approprite with TLS or IP allow list.
- Place this file in the /etc/nginx/conf.d folder, and reload nginx.  The Plus dashboard is now accessible at http://nginx-lb-server-ip:9000/dashboard.html.  It should look similar to this:

![NGINX Dashboard](../media/nkl-http-dashboard.png)

- Use the included nginx.conf file, it enables the NGINX NJS module, for exporting the Plus statistics:  

```bash
cat /etc/nginx/nginx.conf

# NGINX K8s Loadbalancer Solution
# Chris Akker, Apr 2023
# Example nginx.conf
# Enable Prometheus NJS module, increase output buffer size
# Enable Stream context, add /var/log/nginx/stream.log
#
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;

load_module modules/ngx_http_js_module.so;   # Load NJS module

worker_rlimit_nofile 2048;

events {
    worker_connections 2048;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $upstream_addr - $upstream_status - $remote_user [$time_local] $host - "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;

    #added for Prometheus
    subrequest_output_buffer_size 32k;

}

# TCP load balancing block
#
stream {
   include /etc/nginx/stream/*.conf;
    log_format  stream  '$remote_addr - $server_addr [$time_local] $status $upstream_addr $upstream_bytes_sent';
    access_log  /var/log/nginx/stream.log  stream;
}

```

- Configure NGINX for HTTP processing, load balancing, and MultiCluster split clients for this Solution.

  `Notice that this Solution only uses port 443 and terminates TLS.`  
  
  Place the `clusters.conf` file in the /etc/nginx/conf.d folder, and reload NGINX.  Notice the match block and health check directives are for the cafe.example.com Demo application from NGINX.

```bash
cat /etc/nginx/conf.d/clusters.conf

# NGINX K8sLB HTTP configuration, for L7 load balancing
# Chris Akker, Apr 2023
# HTTP Proxy and load balancing
# MultiCluster Load Balancing with http split clients 0-100%
# Upstream servers managed by NKL Controller
# NGINX Key Value store for Split ratios
#
#### clusters.conf

# Define Key Value store, backup state file, timeout, and enable sync

keyval_zone zone=split:1m state=/var/lib/nginx/state/split.keyval timeout=30d sync;
keyval $host $split_level zone=split;

# Main NGINX Server Block for cafe.example.com, with TLS

server {
   listen 443 ssl;
   status_zone https://cafe.example.com;
   server_name cafe.example.com;
   
   ssl_certificate /etc/ssl/nginx/default.crt;  # self-signed for example only
   ssl_certificate_key /etc/ssl/nginx/default.key;
   
   location / {
   status_zone /;
   
   proxy_set_header Host $host;
   proxy_http_version 1.1;
   proxy_set_header   "Connection" "";
   proxy_pass https://$upstream;
   
   }

   location @health_check_cluster1_cafe {

   health_check interval=10 match=cafe;
   proxy_connect_timeout 2s;
   proxy_read_timeout 3s;
   proxy_set_header Host cafe.example.com;
   proxy_pass https://cluster1-https;
   }
 
   location @health_check_cluster2_cafe {

   health_check interval=10 match=cafe;
   proxy_connect_timeout 2s;
   proxy_read_timeout 3s;
   proxy_set_header Host cafe.example.com;
   proxy_pass https://cluster2-https;
   }
}

match cafe {
  status 200-399;
  }

# Cluster1 upstreams

upstream cluster1-https {
   zone cluster1-https 256k;
   least_time last_byte;        # advanced NGINX LB algorithm
   keepalive 16;
   #servers managed by NKL Controller
   state /var/lib/nginx/state/cluster1-https.state; 
}

# Cluster2 upstreams

upstream cluster2-https {
   zone cluster2-https 256k;
   least_time last_byte;        # advanced NGINX LB algorithm
   keepalive 16;
   #servers managed by NKL Controller
   state /var/lib/nginx/state/cluster2-https.state; 
}

# HTTP Split Clients Configuration for Cluster1/Cluster2 ratios
# Ratios provided:  0,1,5,10,25,50,75,90,95,99,100%

split_clients $request_id $split0 { 
   * cluster2-https;
   }

split_clients $request_id $split1 { 
   1.0% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split5 {
   5.0% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split10 { 
   10% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split25 { 
   25% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split50 { 
   50% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split75 { 
   75% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split90 {
   90% cluster1-https;
   * cluster2-https;
   }
   
split_clients $request_id $split95 {
   95% cluster1-https;
   * cluster2-https;
   }
   
split_clients $request_id $split99 {
   99% cluster1-https;
   * cluster2-https;
   }

split_clients $request_id $split100 {
   * cluster1-https;
   }

map $split_level $upstream { 
   0 $split0;
   1.0 $split1;
   5.0 $split5;
   10 $split10;
   25 $split25;
   50 $split50;
   75 $split75;
   90 $split90;
   95 $split95;
   99 $split99;
   100 $split100;
   default $split50;
   }

```

- Configure NGINX for the Prometheus scraper page, which exports the Plus statistics.  Place the `prometheus.conf` file in /etc/nginx/conf.d folder and reload NGINX.

```bash
cat /etc/nginx/conf.d/prometheus.conf

# NGINX K8sLB Prometheus configuration, for HTTP scraper page
# Chris Akker, Apr 2023
# https://www.nginx.com/blog/how-to-visualize-nginx-plus-with-prometheus-and-grafana/
#
#### prometheus.conf

js_import /usr/share/nginx-plus-module-prometheus/prometheus.js;

server {
    location = /metrics {
        js_content prometheus.metrics;
    }

    location /api {
        api;
    } 

}

```

- High Availability:  If you have 2 or more NGINX Plus LB Servers, you can use Zone Sync to synchronize the Split Key Value Store data between the NGINX Servers automatically.  Use the `zonesync.conf` example file provided, change the IP addresses to match your NGINX LB Servers.  Place this file in /etc/nginx/stream folder, and reload NGINX.  Note:  This example does not provide any security for the Zone Sync traffic, secure as necessary with TLS or IP allowlist.

```bash
cat zonesync.conf

# NGINX K8sLB Zone Sync configuration, for KVstore split
# Chris Akker, Apr 2023
# Stream Zone Sync block
# 2 NGINX Plus nodes KVstore zone
# NGINX Kubernetes Loadbalancer
# https://docs.nginx.com/nginx/admin-guide/high-availability/zone_sync/
#
#### zonesync.conf

server {
   zone_sync;

   listen 9001;

   # cluster of 2 nodes
   zone_sync_server 10.1.1.4:9001;
   zone_sync_server 10.1.1.5:9001;

   }

```

Watching the NGINX Plus Dashboard, you will see messages sent/received if Zone Synch is operating correctly:

![Zone Sync](../media/nkl-zone-sync.png)

<br/>

## 5. Install NKL - NGINX Kubernetes Loadbalancer Controller

<br/>

![NIC](../media/nkl-logo.png)

<br/>

This is the new K8s Controller from NGINX, which is configured to watch the k8s environment, the `nginx-ingress Service` object, and send API updates to the NGINX LB Server when there are changes.  It only requires three things.

- New kubernetes namespace and RBAC
- NKL ConfigMap, to configure the Controller
- NKL Deployment, to deploy and run the Controller

Create the new `nkl` K8s namespace:

```bash
kubectl create namespace nkl
```

Apply the manifests for NKL's Secret, Service, ClusterRole, and ClusterRoleBinding:

```bash
kubectl apply -f secret.yaml serviceaccount.yaml clusterrole.yaml clusterrolebinding.yaml
```

Modify the ConfigMap manifest to match your Network environment. Change the `nginx-hosts` IP address to match your NGINX LB Server IP.  If you have 2 or more LB Servers, separate them with a comma.  Important! - keep the port number for the Plus API endpoint, and the `/api` URL as shown.

```yaml
apiVersion: v1
kind: ConfigMap
data:
  nginx-hosts:
    "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api" # change IP(s) to match NGINX LB Server(s)
metadata:
  name: nkl-config
  namespace: nkl
```

Apply the updated ConfigMap:

```bash
kubectl apply -f nkl-configmap.yaml
```

Deploy the NKL Controller:

```bash
kubectl apply -f nkl-deployment.yaml
```

Check to see if the NKL Controller is running with the updated ConfigMap:

```bash
kubectl get pods -n nkl
```
```bash
kubectl describe cm nkl-config -n nkl
```

The status should show "running", your `nginx-hosts` should have the <LB Server IP>:Port/api defined.

![NKL Running](../media/nkl-configmap.png)

To make it easy to watch the NKL Controller log messages, add the following bash alias:

```bash
alias nkl-follow-logs='kubectl -n nkl get pods | grep nkl-deployment | cut -f1 -d" "  | xargs kubectl logs -n nkl --follow $1'
```

Using a new Terminal, watch the NKL Controller log:

```bash
nkl-follow-logs
```

Leave this Terminal window open, so you can watch the log messages!

- Create the NKL compatible NODEPORT Service, using the `nodeport-cluster1.yaml` manifest provided:

```bash
kubectl apply -f nodeport-cluster1.yaml
```

Verify that the `nginx-ingress` NodePort Service is properly defined:

```bash
kubectl get svc nginx-ingress -n nginx-ingress
```

![NGINX Ingress NodePort Service](../media/nkl-cluster1-nodeport.png)
![NGINX Ingress NodePort Service](../media/nkl-cluster1-upstreams.png)

### NodePort is 443:30267, K8s Workers are 10.1.1.8 and .10.

<br/>
<br/>

### MultiCluster Solution

If you plan to implement and test the MultiCluster Load Balancing Solution, repeat all the steps to configure the second K8s cluster, identical to the first Cluster1 steps.  There is only one change - you MUST use the appropriate `nodeport-clusterX.yaml` manifest to match the appropriate cluster.  Don't forget to check and set your ./kube Config Context when you change clusters!

<br/>

## 6. Testing NKL

<br/>

When you are finished, the NGINX Plus Dashboard on the LB Server should look similar to the following image:

![NGINX Upstreams Dashboard](../media/nkl-multicluster-upstreams.png)

Important items for reference:
- Orange are the upstream server blocks, from the `etc/nginx/conf.d/clusters.conf` file.
- If both NKL Controllers are working, it will update the correct `clusterX-https` upstream block. 
- The IP addresses will match the K8s worker nodes, the port numbers will match the NodePort definitions for nginx-ingress Service from each cluster.

>Note: In this example, there is a 3-Node K8s cluster, with one Control Node, and 2 Worker Nodes.  The NKL Controller only configures NGINX with `Worker Node` IP addresses, from Cluster1, which are:
- 10.1.1.8
- 10.1.1.10

Cluster2 Worker Node addresses are:
- 10.1.1.11
- 10.1.1.12

Notice: K8s Control Nodes are excluded from the list intentionally.

<br/>

Configure DNS, or the local hosts file, for cafe.example.com > NGINXLB Server IP Address.  In this example:

```bash
cat /etc/hosts

10.1.1.4 cafe.example.com
```

Open a browser tab to https://cafe.example.com/coffee.  

The Dashboard's `HTTP Upstreams Requests counters` will increase as you refresh the browser page.

Using a Terminal and `./kube Context set for Cluster1`, delete the `nginx-ingress nodeport service` definition.  

```bash
kubectl delete -f nodeport-cluster1.yaml
```

Now the `nginx-ingress` Service is gone, and the Cluster1 upstream list will now be empty in the Dashboard.

![NGINX No Cluster1 NodePort](../media/nkl-cluster1-delete-nodeport.png)
Legend:
- Orange highlights the Cluster1 and NodePort are deleted.
- Indigo highlights the NKL Controller log message, successfully deleting the cluster1-https upstreams.
- Blue highlights the actual API calls to the LB Server, 10.1.1.4.
- Notice there are 4 Delete Log messages, 2 Worker Nodes X 2 LB Servers.
- If you are running a second NGINX LB Server for HA, and Zone Sync is working, the cluster1-https upstreams on LB Server#2 will also be empty.  Check the LB Server#2 Dashboard to confirm.

If you refresh the cafe.example.com browser page, 1/2 of the requests will respond with `502 Bad Gateway`.  There are NO upstreams in Cluster1 for NGINX to send the requests to!

Add the `nginx-ingress` Service back to Cluster1:

```
kubectl apply -f nodeport-cluster1.yaml
```

Verify the nginx-ingress Service is re-created.  Notice the the Port Numbers have changed!

`The NKL Controller detects this change, and modifies the LB Server upstreams.`  The Dashboard will show you the new Port numbers, matching the new NodePort definitions.  The NKL logs show these messages, confirming the changes:

![NKL Add NodePort](../media/nkl-cluster1-add-nodeport.png)

<br/>

## 7. Testing MultiCluster Loadbalancing with HTTP Split Clients

In this section, you will generate some HTTP load on the NGINX LB Server, and watch as it sends traffic to both Clusters.  Then you will `dynamically change the Split ratio`, and watch NGINX send different traffic levels to each cluster.

The only tool you need for this, is an HTTP load generation tool.  WRK, running in a docker container outside the cluster is what is shown here.

Start WRK, on a client outside the cluster.  This command runs WRK for 15 minutes, targets the NGINX LB Server URL of https://10.1.1.4/coffee.  The host header is required, cafe.example.com, as NGINX is configured for this server_name. (And so is the NGINX Ingress Controller).

```bash
docker run --rm williamyeh/wrk -t2 -c200 -d15m -H 'Host: cafe.example.com' --timeout 2s https://10.1.1.4/coffee
```

![nkl Clusters 50-50](../media/nkl-clusters-50.png)

You see the traffic is load balanced between cluster1 and cluster2 at 50/50 ratio.  

Add a record to the KV store, by sending an API command to NGINX Plus:

```bash
curl -iX POST -d '{"cafe.example.com":50}' http://nginxlb:9000/api/8/http/keyvals/split
```

Verify the API record is there, on both NGINX LB Servers:
```bash
curl http://nginxlb:9000/api/8/http/keyvals/split
curl http://nginxlb2:9000/api/8/http/keyvals/split
```

![NGINXLB KeyVal](../media/nkl-keyval-split.png)

If the KV data is missing on one LB Server, your Zone Sync must be fixed.

>Notice the difference in HTTP Response Times, Cluster2 is running much faster than Cluster1 !  (The Red and Green highlights on the Dashboard)

So, you decide to send less traffic to Cluster1, and more to Cluster2.  You will set the HTTP Split ratio to 10/90 = 10% to Cluster1, 90% to Cluster2.

Remember:  This Solution example configures NGINX for Cluster1 to use the Split value, and the remaining percentage of traffic is sent to Cluster2.

Change the KV Split Ratio to 10:
```bash
curl -iX PATCH -d '{"cafe.example.com":10}' http://nginxlb:9000/api/8/http/keyvals/split
```

![NGINXLB Clusters 10-90](../media/nkl-clusters-10.png)

**Important NOTE:**  The first time, an `HTTP POST` is required to ADD a new record to the KeyValue store.  Once the record exists, use an `HTTP PATCH` method to update an existing record, which will change the ratio value in memory, dynamically, with no reloads or restarts of NGINX required!

Try a few more ratios, see how it works.  If you review the `clusters.conf` file, you will discover what Ratios are provided for you.  You can edit to suit your needs, of course.  Notice the Map directive has a "default" set to "50".  So if you make a mistake, it will Split at a 50:50 ratio.

As you can see, if you set the Ratio to "0", `Cluster1 receives NO TRAFFIC`, and you can perform k8s maintenance, troubleshooting, upgrades, etc, with no impact to live traffic.  Alternatively, you can set the Ratio to "100", and now `Cluster2 receives NO TRAFFIC`, and you can work on that cluster - with NO downtime required.

Set the Split back to "50" when your testing is completed, and ask the boss for a raise.

The Completes the Testing Section.

</br>

## 8. Prometheus and Grafana Servers

<br/>

Prometheus | Grafana

![](../media/prometheus-icon.png)  |![](../media/grafana-icon.png)
--- | ---

### During the testing of the Solution, it is helpful to see the load balancing and HTTP Split Ratios visually, using a chart or graph over time.

<br/>

Here are the instructions to run 2 Docker containers on a Monitor Server, which will collect the NGINX Plus statistics from Prometheus, and graph them with Grafana.  Likely, you already have these running in your environment, but are provided here as an example to display the NGINX metrics of high value.

### Prometheus

- Configure your Prometheus server to collect NGINX Plus statistics from the scraper page.  Use the prometheus.yml file provided, edit the IP addresses to match your NGINX LB Server(s).

```bash
cat prometheus.yaml
```

```yaml
global:
  scrape_interval: 15s 
  
  external_labels:
    monitor: 'nginx-monitor'
 
scrape_configs:  
  - job_name: 'prometheus'
    
    scrape_interval: 5s
 
    static_configs:
      - targets: ['10.1.1.4:80', '10.1.1.5:80']  # NGINX LB Servers
```

- Review, edit and place the sample `prometheus.yml` file in /etc/prometheus folder.

- Start the Prometheus docker container:

```bash
sudo docker run --restart always --network="host" -d -p 9090:9090 --name=prometheus -v ~/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

Prometheus Web Console access to the data is on <monitor-server-ip:9090>.

<br/>

### Grafana 

- Create a docker volume to store the Grafana data.

```bash
docker volume create grafana-storage
```

- Start the Grafana docker container:

```bash
sudo docker run --restart always -d -p 3000:3000 --name=grafana -v grafana-storage:/var/lib/grafana grafana/grafana
```

Web console access to Grafana is on <monitor-server-ip:3000>.  Login is admin/admin.

You can import the provided `grafana-dashboard.json` file to see the NGINX Plus `Cluster1 and 2 statistics` HTTP RPS and Upstream Response Times.

<br/>

### For example, here are four Grafana Charts of the MultiCluster Solution, under load with WRK, as detailed above in the Testing Section:

- The Split ratio was started at 50%.  Then it was changed to `10` using the Plus API at 12:56:30.  And then change to `90` 3 minutes later at 12:59:30.
- The first graph is Split Ratio=50, second graph is Split Ratio=10, and third graph is Split Ratio=90.
- The fourth graph is the HTTP Response Time from both Clusters ... why is Cluster2 so much faster than Cluster1 ???  Good thing NGINX provides way to monitor and adjust the traffic based on real-time metrics :-)

![NGINX LB Split 50 Grafana](../media/nkl-grafana-reqs-50.png)
![NGINX LB Split 10 Grafana](../media/nkl-grafana-reqs-10.png)
![NGINX LB Split 90 Grafana](../media/nkl-grafana-reqs-90.png)

![NGINX LB Resp Time Grafana](../media/nkl-grafana-resp.png)

<br/>

### End of Prometheus and Grafana Section

<br/>

## Authors
- Chris Akker - Solutions Architect - Community and Alliances @ F5, Inc.
- Steve Wagner - Solutions Architect - Community and Alliances @ F5, Inc.