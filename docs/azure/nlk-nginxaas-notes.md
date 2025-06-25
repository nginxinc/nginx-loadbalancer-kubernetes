1. Set APIkey and Hostname

export keyValue=123456
echo $keyValue

export dataplaneAPIEndpoint=http://10.1.1.12:9000/
echo $dataplaneAPIEndpoint

2. Helm Install

kubectl create namespace nlk
kubectl get ns -A

export HELM_EXPERIMENTAL_OCI=1
echo $HELM_EXPERIMENTAL_OCI

helm install nlk oci://registry-1.docker.io/nginxcharts/nginxaas-loadbalancer-kubernetes --version 1.0.0 --set "nlk.dataplaneApiKey=${keyValue}" --set "nlk.config.nginxHosts=${dataplaneAPIEndpoint}nplus" --namespace nlk

Optional: Enable NLK Debug Logging

helm install nlk oci://registry-1.docker.io/nginxcharts/nginxaas-loadbalancer-kubernetes --version 1.0.0 --set "nlk.dataplaneApiKey=${keyValue}" --set "nlk.config.nginxHosts=${dataplaneAPIEndpoint}nplus" --set "nlk.config.logLevel=debug" --namespace nlk

helm ls -A

NAME	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART                                 	APP VERSION
nlk 	nlk      	1       	2025-04-01 21:44:35.787452842 +0000 UTC	deployed	nginxaas-loadbalancer-kubernetes-1.0.0	1.0.0

Optional:  NLK updates 2 different VMs, separate nginxHosts with a space character

helm install nlk oci://registry-1.docker.io/nginxcharts/nginxaas-loadbalancer-kubernetes --version 1.0.0 --set "nlk.dataplaneApiKey=123456" --set "nlk.config.nginxHosts=http://10.1.1.12:9000/nplus http://10.1.1.13:9000/nplus" --set "nlk.config.logLevel=debug" --namespace nlk


3. Edit NLK Config Map, enable debug logging

kubectl edit cm nlk-nginxaas-loadbalancer-kubernetes-nlk-config -n nlk

apiVersion: v1
data:
  config.yaml: |
    log-level: "debug"
    nginx-hosts: "http://10.1.1.12/nplus"
...

Optional: Watch NLK controller logs:

kubectl logs deployment/nlk-nginxaas-loadbalancer-kubernetes -n nlk --follow

4. Create new Upstream Block on Nginx VM, and reload

# configuration file /etc/nginx/conf.d/upstreams.conf:
# Cluster2 upstreams

upstream my-service {
   zone my-service 256k;
   least_time last_byte;
   keepalive 16;
   #servers managed by NLK Controller
   state /var/lib/nginx/state/my-service.state; 
}

5. Create new Dashboard.conf on Nginx VM, if not already exists

6. Create a K8s Service File

ubuntu@k8-jumphost:~$ cat nlk-test-svc.yaml 
apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: nginx-ingress
  annotations:
    # Let the controller know to pay attention to this service.
    # If you are connecting multiple controller the value can be used to distinguish them
    nginx.com/nginxaas: nginxaas
spec:
  # expose a port on the nodes
  type: NodePort
  ports:
    - port: 80
      targetPort: http
      protocol: TCP
      # The port name helps connect to NGINXaaS. It must be prefixed with either `http-` or `stream-`
      # and the rest of the name must match the name of an upstream in that context.
      name: http-my-service
  selector:
    app: nginx-ingress

7. Apply the new Service File

kubectl apply -f nlk-test-svc.yaml

Check it:

kubectl get svc -n nginx-ingress

ubuntu@k8-jumphost:~$ k get svc -n nginx-ingress
NAME         TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
my-service   NodePort   10.96.80.141   <none>        80:32695/TCP   11m

8. Check Nginx Plus Dashboard.  HTTP Upstream Tab should have a `my-service` block, with a list of Kubernetes node IPs and Nodeports

kubectl get nodes -o wide
Use "kubectl options" for a list of global command-line options (applies to all commands).
ubuntu@k8-jumphost:~$ k get nodes -o wide
NAME                    STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
k2control.example.com   Ready    control-plane   23h   v1.31.7   10.1.1.8      <none>        Ubuntu 22.04.2 LTS   5.15.0-1035-kvm   containerd://1.7.27
k2node1.example.com     Ready    <none>          23h   v1.31.7   10.1.1.10     <none>        Ubuntu 22.04.2 LTS   5.15.0-1035-kvm   containerd://1.7.27
k2node2.example.com     Ready    <none>          23h   v1.31.7   10.1.1.11     <none>        Ubuntu 22.04.2 LTS   5.15.0-1035-kvm   containerd://1.7.27

kubectl get svc -n nginx-ingress

ubuntu@k8-jumphost:~$ k get svc -n nginx-ingress
NAME         TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
my-service   NodePort   10.96.80.141   <none>        80:32695/TCP   18m

NOTE:  Service type LoadBalancer is not supported in the Service manifest

###  Monitor Server for Prometheus and Grafana - Build notes

1. Install Docker on ubuntu server

https://docs.docker.com/engine/install/ubuntu/

1. Edit prometheus.yml file

root@ubuntu:/# cat /etc/prometheus/prometheus.yml 

```yaml
global:
  scrape_interval: 15s 
  
  external_labels:
    monitor: 'nginx-monitor'

scrape_configs:  
  - job_name: 'prometheus'
    
    scrape_interval: 5s

    static_configs:
      - targets: ['10.1.1.12:9113', '10.1.1.13:9113']  # NGINX Loadbalancing Servers
```

1. Run Prometheus container

sudo docker run --restart always --network="host" -d -p 9090:9090 --name=prometheus -v /etc/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

1. Run Grafana container



