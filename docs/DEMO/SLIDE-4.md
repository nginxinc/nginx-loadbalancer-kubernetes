## One-way TLS and Mutual TLS

There are two types of TLS: one-way TLS and mutual TLS.

### One-way TLS

One-way TLS is the most common type of TLS. In one-way TLS, the client verifies the server's identity, but the server does not verify the client's identity. One-way TLS is used to secure the connection between the client and the server.

### Mutual TLS

Mutual TLS is less common than one-way TLS. In mutual TLS, the client verifies the server's identity, and the server verifies the client's identity. Mutual TLS is used to secure the connection between the client and the server.

The following diagram shows the difference between one-way TLS and mutual TLS.

```mermaid
graph LR
CACertificate[CA Certificate]

    subgraph "One-way TLS"
        NGINXPlusCert[NGINX Plus Certificate]
        NLK[nginx-loadbalancer-kubernetes]
        NGINXPlus[NGINX Plus]
        NGINXPlusCert -->|Used by| NLK
        NLK -->|Verifies| NGINXPlus
    end

    subgraph "Mutual TLS"
        NLKCert[NLK Certificate]
        MNGINXPlusCert[NGINX Plus Certificate]
        MNLK[nginx-loadbalancer-kubernetes]
        MNGINXPlus[NGINX Plus]
        NLKCert -->|Used by| MNGINXPlus
        MNGINXPlus -->|Verifies| MNLK
        MNGINXPlusCert -->|Used by| MNLK
        MNLK -->|Verifies| MNGINXPlus
    end

CACertificate -->|Used for Signing| NLKCert
CACertificate -->|Used for Signing| NGINXPlusCert
CACertificate -->|Used for Signing| MNGINXPlusCert
```

[Next](SLIDE-5.md)