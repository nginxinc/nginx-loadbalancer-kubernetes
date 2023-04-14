<div>
    <span>
<img src="nkl-logo.svg" width="124" />
</span>    
<span>
<h2>nginx-k8s-loadbalancer</h2>
</span>
</div>

The NGINX K8s Loadbalancer, or _NKL_, is a Kubernetes controller that provides TCP load balancing external to a Kubernetes cluster running on-premise.

## Requirements

### Who needs NKL?

- [ ] If you find yourself living in a world where Kubernetes is running on-premise instead of a cloud provider, you might need NKL.
- [ ] If you want exceptional, best-in-class load-balancing for your Kubernetes applications, you might need NKL.
- [ ] If you want the ability to manage your load-balancing configuration with the same tools you use to manage your Kubernetes cluster, you might need NKL.

### Why NKL?

NKL provides a simple, easy-to-manage way to automate load balancing for your Kubernetes applications by leveraging NGINX Plus hosts running outside your cluster.

NKL installs easily, has a small footprint, and is easy to configure and manage.

? {{review for embetterment}}: NKL does not require any specific domain knowledge for configuration, though you will have to understand NGINX configuration to get the most out of this solution. There is thorough documentation available about these specifics in the `docs/` directory.

### What does NKL do?

tl;dr:

_**NKL is a Kubernetes controller that monitors Services and Nodes in your cluster, and then sends API calls to an external NGINX Plus server to manage NGINX Plus Upstream servers automatically.**_

That's all well and good, but what does that mean? Well, Kubernetes clusters require some tooling to handling routing traffic from the outside world (e.g.: the Internet, corporate network, etc.) to the cluster. 
This is typically done with a load balancer. The load balancer is responsible for routing traffic to the appropriate Kubernetes worker node which then forwards the traffic to the appropriate Service / Pod.

If you are using a hosted web solution -- Digital Ocean, AWS, Azure, etc. -- you can use the cloud provider's load balancer service. This service will create a load balancer for you. 
You can use the cloud provider's API to manage the load balancer, or you can use the cloud provider's web console.

However, if you checked the first box above, you are running Kubernetes on-premise and will need to manage your own load balancer. This is where NKL comes in.

NKL itself does not perform load balancing. Instead, NKL allows you to manage resources within your cluster and have the load balancers automatically be updated to support those changes, with tooling you are most likely already using. 

## Getting Started

There are few bits of administrivia to get out of the way before you can start leveraging NKL for your load balancing needs.

As noted above, NKL really shines when you have one or more Kubernetes clusters running on-premise. With this in place,
you need to have at least one NGINX Plus host running outside your cluster (Please refer to the [Roadmap](#Roadmap) for information about other load balancer servers). 

You will not need to clone this repo to use NKL. Instead, you can install NKL using the included Manifest files (just copy the `deployments/` directory), which pulls the NKL image from the Container Registry.

### RBAC

As with everything Kubernetes, NKL requires RBAC permissions to function properly. The necessary resources are defined in the various YAML files in `deployement/rbac/`.

For convenience, two scripts are included, `apply.sh`, and `unapply.sh`. These scripts will apply or remove the RBAC resources, respectively.

The permissions required by NKL are modest. NKL requires the ability to read Resources via shared informers; the resources are Services, Nodes, and ConfigMaps. 
The Services and ConfigMap are restricted to a specific namespace (default: "nkl"). The Nodes resource is cluster-wide.

### Configuration

NKL is configured via a ConfigMap, the default settings are found in `deployment/configmap.yaml`. Presently there is a single configuration value exposed in the ConfigMap, `nginx-hosts`. 
This contains a comma-separated list of NGINX Plus hosts that NKL will maintain.

You will need to update this ConfigMap to reflect the NGINX Plus hosts you wish to manage.

If you were to deploy the ConfigMap and start NKL without updating the `nginx-hosts` value, don't fear; the ConfigMap resource is monitored for changes and NKL will update the NGINX Plus hosts accordingly when the resource is changed, no restart required.

### Deployment

There is an extensive [Installation Guide](docs/InstallationGuide.md) available in the `docs/` directory. 
Please refer to that for detailed instructions on how to deploy NKL and run a demo application.

To get NKL up and running in ten steps or fewer, follow these instructions (NOTE, all the aforementioned prerequisites must be met for this to work):

1. Clone this repo (optional, you can simply copy the `deployments/` directory) 

```git clone git@github.com:nginxinc/nginx-k8s-loadbalancer.git```

2. Apply the RBAC resources

```./deployments/rbac/apply.sh```

3. Apply the Namespace

```kubectl apply -f deployments/namespace.yaml```

4. Update / Apply the ConfigMap (For best results update the `nginx-hosts` values first)

```kubectl apply -f deployments/configmap.yaml```

5. Apply the Deployment

```kubectl apply -f deployments/deployment.yaml```

6. Check the logs

```kubectl -n nkl get pods | grep nkl-deployment | cut -f1 -d" "  | xargs kubectl logs -n nkl --follow $1```

At this point NKL should be up and running. Now would be a great time to go over to the [Installation Guide](docs/InstallationGuide.md) 
and follow the instructions to deploy a demo application.

### Monitoring

Presently NKL includes a fair amount of logging. This is intended to be used for debugging purposes. 
There are plans to add more robust monitoring and alerting in the future.

As a rule, we support the use of [OpenTelemetry](https://opentelemetry.io/) for observability, and we will be adding support in the near future.

## Contributing

Presently we are not accepting pull requests. However, we welcome your feedback and suggestions. 
Please open an issue to let us know what you think!

## Roadmap

While NKL was initially written specifically for NGINX Plus, we recognize there are other load-balancers that can be supported.

To this end, NKL has been architected to be extensible to support other "Border Servers". 
Border Servers are the term NKL uses to describe load-balancers, reverse proxies, etc. that run outside the cluster and handle 
routing outside traffic to your cluster. 

While we have identified a few potential targets, we are open to suggestions. Please open an issue to share your thoughts on potential targets.

We look forward to building a community around NKL and value all feedback and suggestions. Varying perspectives and embracing
diverse ideas will be key to NKL becoming a solution that is useful to the community. We will consider it a success
when we are able to accept pull requests from the community.

## License

[Apache License, Version 2.0](https://github.com/nginxinc/nginx-k8s-loadbalancer/blob/main/LICENSE)

&copy; [F5, Inc.](https://www.f5.com/) 2023

(but don't let that scare you, we're really nice people...)
