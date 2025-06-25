[![CI](https://github.com/nginxinc/nginx-loadbalancer-kubernetes/actions/workflows/build-test.yml/badge.svg)](https://github.com/nginxinc/nginx-loadbalancer-kubernetess/actions/workflows/build-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nginxinc/nginx-loadbalancer-kubernetes)](https://goreportcard.com/report/github.com/nginxinc/nginx-loadbalancer-kubernetes)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nginxinc/nginx-loadbalancer-kubernetes?logo=github&sort=semver)](https://github.com/nginxinc/nginx-loadbalancer-kubernetes/releases/latest) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nginxinc/nginx-loadbalancer-kubernetes?logo=go)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nginxinc/nginx-loadbalancer-kubernetes/badge)](https://api.securityscorecards.dev/projects/github.com/nginxinc/nginx-loadbalancer-kubernetes)
[![CodeQL](https://github.com/nginxinc/nginx-loadbalancer-kubernetes/workflows/codeql.yml/badge.svg?branch=main&event=push)](https://github.com/nginxinc/nginx-loadbalancer-kubernetes/actions/codeql.yml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fnginxinc%2Fnginx-loadbalancer-kubernetes.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fnginxinc%2Fnginx-loadbalancer-kubernetes?ref=badge_shield)
[![Community Support](https://badgen.net/badge/support/community/cyan?icon=awesome)](https://github.com/nginxinc/nginx-loadbalancer-kubernetes/discussions)
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)



<div style="margin-bottom: 5em;">
    <span>
        <img style="float: left;" src="nlk-logo.svg" width="124" />
        <h2 style="padding: 1.5em">nginx-loadbalancer-kubernetes</h2>
    </span>
</div>


The NGINX Loadbalancer for Kubernetes, or _NLK_, is a Kubernetes controller that provides TCP load balancing external to a Kubernetes cluster running on-premise.

## Requirements

[//]: # (### Who needs NLK?)

[//]: # ()
[//]: # (- [ ] If you find yourself living in a world where Kubernetes is running on-premise instead of a cloud provider, you might need NLK.)

[//]: # (- [ ] If you want exceptional, best-in-class load-balancing for your Kubernetes clusters by using NGINX Plus, you might need NLK.)

[//]: # (- [ ] If you want the ability to manage your load-balancing configuration with the same tools you use to manage your Kubernetes cluster, you might need NLK.)

### What you will need

- [ ] A Kubernetes cluster running on-premise.
- [ ] One or more NGINX Plus hosts running outside your Kubernetes cluster (NGINX Plus hosts must have the ability to route traffic to the cluster).

There is a more detailed [Installation Reference](docs/README.md) available in the `docs/` directory.

### Why NLK?

NLK provides a simple, easy-to-manage way to automate load balancing for your Kubernetes applications by leveraging NGINX Plus hosts running outside your cluster.

NLK installs easily, has a small footprint, and is easy to configure and manage.

NLK does not require learning a custom object model, you only have to understand NGINX configuration to get the most out of this solution.
There is thorough documentation available with the specifics in the `docs/` directory.

### What does NLK do?

tl;dr:

_**NLK is a Kubernetes controller that monitors Services and Nodes in your cluster, and then sends API calls to an external NGINX Plus server to manage NGINX Plus Upstream servers automatically.**_

That's all well and good, but what does it mean? Kubernetes clusters require some tooling to handling routing traffic from the outside world (e.g.: the Internet, corporate network, etc.) to the cluster.
This is typically done with a load balancer. The load balancer is responsible for routing traffic to the appropriate worker node which then forwards the traffic to the appropriate Service / Pod.

If you are using a hosted Kubernetes solution -- Digital Ocean, AWS, Azure, etc. -- you can use the cloud provider's load balancer service. Those services will create a load balancer for you.
You can use the cloud provider's API to manage the load balancer, or you can use the cloud provider's web console.

If you are running Kubernetes on-premise and will need to manage your own load balancer, NLK can help.

NLK itself does not perform load balancing. Rather, NLK allows you to manage Service resources within your cluster to update your load balancers, with tooling you are most likely already using.

<img src="docs/media/nlk-blog-diagram-v1.png" width="768" />

## Getting Started

There are few bits of administrivia to get out of the way before you can start leveraging NLK for your load balancing needs.

As noted above, NLK is intended for when you have one or more Kubernetes clusters running on-premise. In addition to this,
you need to have at least one NGINX Plus host running outside your cluster (Please refer to the [Roadmap](#Roadmap) for information about other load balancer servers).

### Deployment

#### RBAC

As with everything Kubernetes, NLK requires RBAC permissions to function properly. The necessary resources are defined in the various YAML files in `deployment/rbac/`.

For convenience, two scripts are included, `apply.sh`, and `unapply.sh`. These scripts will apply or remove the RBAC resources, respectively.

The permissions required by NLK are modest. NLK requires the ability to read Resources via shared informers; the resources are Services, Nodes, and ConfigMaps.
The Services and ConfigMap are restricted to a specific namespace (default: "nlk"). The Nodes resource is cluster-wide.

#### Configuration

NLK is configured via a ConfigMap, the default settings are found in `deployment/configmap.yaml`. Presently there is a single configuration value exposed in the ConfigMap, `nginx-hosts`.
This contains a comma-separated list of NGINX Plus hosts that NLK will maintain.

You will need to update this ConfigMap to reflect the NGINX Plus hosts you wish to manage.

If you were to deploy the ConfigMap and start NLK without updating the `nginx-hosts` value, don't fear; the ConfigMap resource is monitored for changes and NLK will update the NGINX Plus hosts accordingly when the resource is changed, no restart required.

There is an extensive [Installation Reference](docs/README.md) available in the `docs/` directory.
Please refer to that for detailed instructions on how to deploy NLK and run a demo application.

#### Versioning

Versioning is a work in progress. The CI/CD pipeline is being developed and will be used to build and publish NLK images to the Container Registry.
Once in place, semantic versioning will be used for published images.

#### Deployment Steps

To get NLK up and running in ten steps or fewer, follow these instructions (NOTE, all the aforementioned prerequisites must be met for this to work).
There is a much more detailed [Installation Reference](docs/README.md) available in the `docs/` directory.

1. Clone this repo (optional, you can simply copy the `deployments/` directory)

```git clone git@github.com:nginxinc/nginx-loadbalancer-kubernetes.git```

2. Apply the Namespace

```kubectl apply -f deployments/deployment/namespace.yaml```

3. Apply the RBAC resources

```./deployments/rbac/apply.sh```

4. Update / Apply the ConfigMap (For best results update the `nginx-hosts` values first)

```kubectl apply -f deployments/deployment/configmap.yaml```

5. Apply the Deployment

```kubectl apply -f deployments/deployment/deployment.yaml```

6. Check the logs

```kubectl -n nlk logs -f $(kubectl -n nlk get po -l "app=nlk" --no-headers -o custom-columns=":metadata.name")```

At this point NLK should be up and running. Now would be a great time to go over to the [Installation Reference](docs/README.md)
and follow the instructions to deploy a demo application.

### Monitoring

Presently NLK includes a fair amount of logging. This is intended to be used for debugging purposes.
There are plans to add more robust monitoring and alerting in the future.

As a rule, we support the use of [OpenTelemetry](https://opentelemetry.io/) for observability, and we will be adding support in the near future.

## Contributing

Presently we are not accepting pull requests. However, we welcome your feedback and suggestions.
Please open an issue to let us know what you think!

One way to contribute is to help us test NLK. We are looking for people to test NLK in a variety of environments.

If you are curious about the implementation, you should certainly browse the code, but first you might wish to refer to the [design document](docs/DESIGN.md).
Some of the design decisions are explained there.

## Roadmap

While NLK was initially written specifically for NGINX Plus, we recognize there are other load-balancers that can be supported.

To this end, NLK has been architected to be extensible to support other "Border Servers".
Border Servers are the term NLK uses to describe load-balancers, reverse proxies, etc. that run outside the cluster and handle
routing outside traffic to your cluster.

While we have identified a few potential targets, we are open to suggestions. Please open an issue to share your thoughts on potential implementations.

We look forward to building a community around NLK and value all feedback and suggestions. Varying perspectives and embracing
diverse ideas will be key to NLK becoming a solution that is useful to the community. We will consider it a success
when we are able to accept pull requests from the community.

## Authors
- Chris Akker - Solutions Architect - Community and Alliances @ F5, Inc.
- Steve Wagner - Solutions Architect - Community and Alliances @ F5, Inc.

<br/>

## License

[Apache License, Version 2.0](https://github.com/nginxinc/nginx-loadbalancer-kubernetes/blob/main/LICENSE)

&copy; [F5, Inc.](https://www.f5.com/) 2023

(but don't let that scare you, we're really nice people...)
