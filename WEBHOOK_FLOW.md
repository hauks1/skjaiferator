# Conversion Webhook Flow Diagrams

## 1. Runtime Flow: v1alpha1 → v1beta1 Conversion

```mermaid
sequenceDiagram
    participant User
    participant K8s API Server
    participant Webhook Service
    participant Operator Pod
    participant etcd
    participant Reconciler

    User->>K8s API Server: kubectl apply v1alpha1 resource
    K8s API Server->>K8s API Server: Check CRD definition<br/>conversion.strategy = Webhook<br/>storage version = v1beta1
    K8s API Server->>Webhook Service: POST /convert<br/>desiredAPIVersion: v1beta1<br/>objects: [v1alpha1]
    Webhook Service->>Operator Pod: Route to port 9443
    Operator Pod->>Operator Pod: controller-runtime<br/>calls ConvertTo()<br/>v1alpha1 → v1beta1
    Operator Pod-->>Webhook Service: Return v1beta1 object
    Webhook Service-->>K8s API Server: convertedObjects: [v1beta1]
    K8s API Server->>etcd: Store v1beta1<br/>(storage version)
    etcd-->>K8s API Server: Stored
    K8s API Server->>Reconciler: Watch event<br/>(v1beta1 from cache)
    Reconciler->>Reconciler: Process v1beta1 object
```

## 2. Kustomize Configuration Chain

```mermaid
graph TD
    A[config/default/kustomization.yaml] --> B[config/crd]
    A --> C[config/webhook]
    A --> D[config/certmanager]
    A --> E[config/manager]
    A --> F[patches: manager_webhook_patch.yaml]
    A --> G[replacements: CA injection]

    B --> B1[bases/svartskjaifs.yaml<br/>Generated CRD]
    B --> B2[patches/webhook_in_svartskjaifs.yaml<br/>Adds conversion.webhook config]
    
    C --> C1[service.yaml<br/>webhook-service<br/>port 443 → 9443]
    
    D --> D1[certificate-webhook.yaml<br/>serving-cert<br/>dnsNames: webhook-service.svc]
    
    E --> E1[manager.yaml<br/>operator deployment]
    
    F --> F1[Adds webhook server port 9443<br/>to manager deployment]
    
    G --> G1[Certificate namespace/name<br/>→ CRD annotation<br/>cert-manager.io/inject-ca-from]
    
    style A fill:#e1f5ff
    style B2 fill:#ffe1e1
    style C1 fill:#e1ffe1
    style D1 fill:#fff4e1
    style G1 fill:#f0e1ff
```

## 3. CRD Conversion Configuration Flow

```mermaid
graph LR
    A[webhook_in_svartskjaifs.yaml] -->|patches| B[CRD]
    
    B --> B1[spec.conversion.strategy: Webhook]
    B --> B2[spec.conversion.webhook.clientConfig]
    
    B2 --> C1[service.name: webhook-service]
    B2 --> C2[service.namespace: skjaiferator-system]
    B2 --> C3[service.path: /convert]
    B2 --> C4[caBundle: injected by cert-manager]
    
    D[Certificate: serving-cert] -->|cert-manager<br/>ca-injector| C4
    
    E[replacements in kustomization] -->|creates annotation| F[CRD annotation:<br/>cert-manager.io/inject-ca-from]
    
    F -->|watched by| D
    
    style B fill:#e1f5ff
    style C4 fill:#ffe1e1
    style D fill:#fff4e1
```

## 4. Cert-Manager CA Injection Flow

```mermaid
sequenceDiagram
    participant Kustomize
    participant K8s API
    participant cert-manager
    participant ca-injector
    participant CRD

    Kustomize->>K8s API: Apply resources with replacements
    K8s API->>K8s API: Create Certificate: serving-cert
    K8s API->>K8s API: Create CRD with annotation:<br/>cert-manager.io/inject-ca-from:<br/>skjaiferator-system/serving-cert
    
    cert-manager->>cert-manager: Generate self-signed cert
    cert-manager->>K8s API: Store in Secret: webhook-server-cert
    
    ca-injector->>CRD: Watch CRD with inject-ca-from annotation
    ca-injector->>K8s API: Read Secret: webhook-server-cert
    ca-injector->>CRD: Inject CA bundle into<br/>spec.conversion.webhook.clientConfig.caBundle
    
    Note over CRD: Now K8s API can verify<br/>webhook server's TLS cert
```

## 5. Component Connections

```mermaid
graph TB
    subgraph "Kubernetes Cluster"
        A[API Server]
        B[etcd]
        
        subgraph "skjaiferator-system namespace"
            C[webhook-service<br/>Service]
            D[controller-manager<br/>Pod]
            E[webhook-server-cert<br/>Secret]
        end
        
        subgraph "cert-manager namespace"
            F[cert-manager<br/>Pod]
            G[ca-injector<br/>Pod]
        end
    end
    
    H[CRD:<br/>svartskjaifs.skjaif.skjaiferator.no]
    I[Certificate:<br/>serving-cert]
    
    A -->|calls when converting| C
    C -->|routes to port 9443| D
    D -->|reads TLS cert from| E
    
    I -->|creates| E
    F -->|manages| I
    
    H -->|has annotation| I
    G -->|injects CA into| H
    
    A -->|reads caBundle from| H
    A -->|stores v1beta1 in| B
    
    style D fill:#e1f5ff
    style C fill:#e1ffe1
    style H fill:#ffe1e1
    style E fill:#fff4e1
```

## 6. Code Registration in main.go

```mermaid
graph TD
    A[main.go] --> B[Register Schemes]
    A --> C[Create Manager]
    A --> D[Setup Controller]
    A --> E[Setup Webhook]
    
    B --> B1[v1alpha1.AddToScheme]
    B --> B2[v1beta1.AddToScheme]
    
    C --> C1[WebhookServer config<br/>port: 9443<br/>cert dir: /tmp/k8s-webhook-server]
    
    D --> D1[SvartSkjaifReconciler<br/>Watches: v1alpha1]
    
    E --> E1[webhookv1beta1.SetupWebhookWithManager]
    E1 --> E2[ctrl.NewWebhookManagedBy<br/>For: v1beta1.SvartSkjaif]
    
    E2 --> E3[v1beta1: Hub marker]
    E2 --> E4[v1alpha1: ConvertTo/ConvertFrom]
    
    style E fill:#e1f5ff
    style E3 fill:#e1ffe1
    style E4 fill:#ffe1e1
```

## Key Files and Their Roles

| File | Purpose |
|------|---------|
| `api/v1alpha1/svartskjaif_conversion.go` | Implements `ConvertTo()` and `ConvertFrom()` |
| `api/v1beta1/svartskjaif_conversion.go` | Implements `Hub()` marker (storage version) |
| `config/crd/patches/webhook_in_svartskjaifs.yaml` | Adds conversion webhook config to CRD |
| `config/webhook/service.yaml` | Creates webhook Service (routes to operator) |
| `config/certmanager/certificate-webhook.yaml` | Requests TLS cert from cert-manager |
| `config/default/kustomization.yaml` | Orchestrates everything with replacements |
| `internal/webhook/v1beta1/svartskjaif_webhook.go` | Registers webhook with controller-runtime |
| `cmd/main.go` | Wires everything together at runtime |

