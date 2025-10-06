# Container Security (CS) Helm Chart

This Helm chart deploys the Container Security platform, a comprehensive runtime monitoring and security solution for Kubernetes environments.

## Introduction

Container Security is a runtime security platform that provides:

- **Runtime Monitoring**: Real-time container activity monitoring using eBPF-based Tetragon
- **Policy Enforcement**: Dynamic security policy management and enforcement
- **Event Processing**: Collection and analysis of security events
- **History and Auditing**: Complete audit trail of security events
- **Multi-cluster Management**: Centralized management of multiple Kubernetes clusters
- **API Access**: RESTful API for integration with external systems

## Architecture

The chart deploys the following components:

### Core Components

- **auth-center**: Authentication and authorization service
- **policy-enforcer**: Security policy enforcement engine
- **history-api**: Historical event query API
- **cluster-manager**: Multi-cluster orchestration
- **notifier**: Alert and notification system
- **reverse-proxy**: Ingress and routing layer
- **cs-manager**: Central management service
- **public-api**: External API gateway

### Runtime Components

- **runtime-monitor**: eBPF-based runtime monitoring agent (Tetragon)
- **event-processor**: Real-time event processing pipeline

### Infrastructure Components (Optional)

- **postgresql**: Primary data store
- **redis**: Cache and session store
- **rabbitmq**: Message broker for event streaming
- **clickhouse**: Analytics database for historical events

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure (if persistence is enabled)
- Sufficient resources for running all components

## Installation

### Quick Start

Install the chart with default configuration:

```bash
helm install cs -n cs ./install/helm \
  --set global.keys.encryption=<YOUR_ENCRYPTION_KEY> \
  --set global.keys.publicAccessTokenSalt=<YOUR_SALT> \
  --set global.ownCsUrl=https://your-domain.com \
  --set auth-center.administrator.username=admin \
  --set auth-center.administrator.password=<YOUR_PASSWORD> \
  --create-namespace
```

**Important**: Replace placeholder values with secure, randomly generated strings. Do not use the example values in production.

### Installing with External Databases

To use external databases instead of deploying them in the cluster:

```bash
helm install cs -n cs ./install/helm \
  --set postgresql.deploy=false \
  --set postgresql.externalHost=postgres.example.com \
  --set redis.deploy=false \
  --set redis.externalHost=redis.example.com \
  --set rabbitmq.deploy=false \
  --set rabbitmq.externalHost=rabbitmq.example.com \
  --set clickhouse.deploy=false \
  --set clickhouse.externalHost=clickhouse.example.com \
  --set global.keys.encryption=<YOUR_ENCRYPTION_KEY> \
  --set global.keys.publicAccessTokenSalt=<YOUR_SALT> \
  --set global.ownCsUrl=https://your-domain.com \
  --set auth-center.administrator.username=admin \
  --set auth-center.administrator.password=<YOUR_PASSWORD> \
  --create-namespace
```

### Installing with Custom Values File

Create a `custom-values.yaml` file:

```yaml
global:
  keys:
    encryption: "<YOUR_ENCRYPTION_KEY>"
    publicAccessTokenSalt: "<YOUR_SALT>"
  ownCsUrl: "https://your-domain.com"

auth-center:
  administrator:
    username: admin
    password: "<YOUR_PASSWORD>"
  replicas: 2

reverse-proxy:
  ingress:
    enabled: true
    class: nginx
    hostname: cs.your-domain.com

postgresql:
  persistence:
    size: 10Gi
    storageClass: fast-ssd

redis:
  persistence:
    size: 5Gi

rabbitmq:
  persistence:
    size: 10Gi

clickhouse:
  persistence:
    size: 50Gi
```

Install with the custom values:

```bash
helm install cs -n cs ./install/helm -f custom-values.yaml --create-namespace
```

## Upgrading

To upgrade an existing installation:

```bash
helm upgrade cs -n cs ./install/helm -f custom-values.yaml
```

## Uninstalling

To uninstall/delete the `cs` deployment:

```bash
helm uninstall cs -n cs
```

This command removes all the Kubernetes components associated with the chart and deletes the release.

**Note**: Persistent Volume Claims are not deleted automatically. To delete them:

```bash
kubectl delete pvc -n cs -l app.kubernetes.io/instance=cs
```

## Configuration

### Security Keys

The following keys must be set for the system to function:

- `global.keys.encryption`: Encryption key for sensitive data in the database (minimum 32 characters)
- `global.keys.token`: Encryption key for authentication tokens (minimum 32 characters)
- `global.keys.publicAccessTokenSalt`: Salt for public API tokens (minimum 32 characters)

#### Auto-generated Keys

You can let Helm generate these keys automatically by setting the value to `INIT-DO-NOT-USE`:

```sh
helm install cs -n cs ./install/helm \
  --set global.keys.encryption=INIT-DO-NOT-USE \
  --set global.keys.publicAccessTokenSalt=INIT-DO-NOT-USE \
  --set global.ownCsUrl=https://your-domain.com \
  --set auth-center.administrator.username=admin \
  --set auth-center.administrator.password=<YOUR_PASSWORD> \
  --create-namespace
```

**Warning**: This approach should be used with extreme caution:
- Keys are generated during installation and stored in Kubernetes secrets
- On subsequent `helm upgrade` operations, if you continue to use `INIT-DO-NOT-USE`, new keys may be regenerated
- Regenerated keys will make existing encrypted data **unrecoverable**, resulting in **data loss**
- This feature is intended for development and testing environments only
- For production environments, always generate and manage keys securely outside of Helm

**Important**: These keys should be securely generated and stored. Changing them after initial deployment may result in data loss.

### TLS Configuration

The chart supports TLS configuration in multiple ways:

#### Auto-generated Self-signed Certificates

```yaml
tls:
  autoGenerated: true
```

#### Existing Secret

```yaml
tls:
  autoGenerated: false
global:
  tls:
    existingSecret: "cs-tls-secret"
```

The secret should contain:
- `tls.crt`: TLS certificate
- `tls.key`: TLS private key
- `ca.crt`: CA certificate (optional)

#### Inline Certificates

```yaml
tls:
  autoGenerated: false
  cert: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
  certKey: |
    -----BEGIN PRIVATE KEY-----
    ...
    -----END PRIVATE KEY-----
  certCA: |
    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
```

### Ingress Configuration

To expose the Container Security web interface:

```yaml
reverse-proxy:
  ingress:
    enabled: true
    class: nginx
    hostname: cs.your-domain.com
    tls:
      autoGenerated: false
      existingSecret: "cs-ingress-tls"
```

### Multi-cluster Setup

#### Primary Cluster

```yaml
global:
  ownCsUrl: "https://cs-primary.your-domain.com"
  isChildCluster: false
```

#### Child Cluster

```yaml
global:
  ownCsUrl: "https://cs-child.your-domain.com"
  centralCsUrl: "https://cs-primary.your-domain.com"
  isChildCluster: true
```

### Resource Requirements

Default resource limits can be adjusted based on your workload:

```yaml
runtime-monitor:
  resources:
    limits:
      cpu: 2
      memory: 4Gi
    requests:
      cpu: 200m
      memory: 256Mi

event-processor:
  resources:
    limits:
      cpu: 2
      memory: 4Gi
    requests:
      cpu: 200m
      memory: 256Mi
```

### High Availability

For production deployments, increase replica counts:

```yaml
auth-center:
  replicas: 3

policy-enforcer:
  replicas: 3

history-api:
  replicas: 3

event-processor:
  replicas: 3

public-api:
  replicas: 3
```

### Monitoring and Metrics

Enable metrics collection:

```yaml
metrics:
  enabled: true

postgresql:
  metrics:
    enabled: true

rabbitmq:
  metrics:
    enabled: true

clickhouse:
  metrics:
    enabled: true
```

### Node Affinity

Deploy components to specific nodes:

```yaml
runtime-monitor:
  nodeSelector:
    workload: monitoring

postgresql:
  nodeSelector:
    workload: database

clickhouse:
  nodeSelector:
    workload: analytics
```

## Parameters

### Global parameters

| Name                                   | Description                                                                                | Value         |
| -------------------------------------- | ------------------------------------------------------------------------------------------ | ------------- |
| `global.imageRegistry`                 | Global Docker image registry                                                               | `""`          |
| `global.imageTag`                      | Global Docker image tag to use for CS components                                           | `""`          |
| `global.logLevel`                      | Logging level for components                                                               | `INFO`        |
| `global.tls.existingSecret`            | Name of the existing secret with TLS certificates                                          | `""`          |
| `global.keys.existingSecret`           | Existing secret name with keys `encryption` and `token`                                    | `""`          |
| `global.keys.encryption`               | Encryption key for secrets stored in database                                              | `""`          |
| `global.keys.token`                    | Encryption key for authentication tokens                                                   | `""`          |
| `global.keys.publicAccessTokenSalt`    | Encryption key for public-api salt                                                         | `""`          |
| `global.postgresql.tls.enabled`        | Enable TLS traffic support (overrides `tls.enabled`)                                       | `true`        |
| `global.postgresql.tls.verify`         | Verify TLS connection to the service (overrides `tls.verify`)                              | `true`        |
| `global.postgresql.tls.existingSecret` | Name of an existing secret that contains the certificates (overrides `tls.existingSecret`) | `""`          |
| `global.redis.tls.enabled`             | Enable TLS traffic support (overrides `tls.enabled`)                                       | `true`        |
| `global.redis.tls.verify`              | Verify TLS connection to the service (overrides `tls.verify`)                              | `true`        |
| `global.redis.tls.existingSecret`      | Name of an existing secret that contains the certificates (overrides `tls.existingSecret`) | `""`          |
| `global.clickhouse.tls.enabled`        | Enable TLS traffic support (overrides `tls.enabled`)                                       | `true`        |
| `global.clickhouse.tls.verify`         | Verify TLS connection to the service (overrides `tls.verify`)                              | `true`        |
| `global.clickhouse.tls.existingSecret` | Name of an existing secret that contains the certificates (overrides `tls.existingSecret`) | `""`          |
| `global.grafana.tls.enabled`           | Enable TLS traffic support (overrides `tls.enabled`)                                       | `true`        |
| `global.grafana.tls.verify`            | Verify TLS connection to the service (overrides `tls.verify`)                              | `true`        |
| `global.grafana.tls.existingSecret`    | Name of an existing secret that contains the certificates (overrides `tls.existingSecret`) | `""`          |
| `global.imagePullSecrets`              | Names of the secrets of the global container registry as an array                          | `["regcred"]` |
| `global.ownCsUrl`                      | URL of primary installation                                                                | `""`          |
| `global.centralCsUrl`                  | URL of primary installation                                                                | `""`          |
| `global.isChildCluster`                | Is this a child cluster                                                                    | `false`       |

### Common CS parameters

| Name                       | Description                                            | Value     |
| -------------------------- | ------------------------------------------------------ | --------- |
| `fullnameOverride`         | String to fully override common.fullname               | `cs`      |
| `imagePullSecret.name`     | Name of the secret with container registry credentials | `regcred` |
| `imagePullSecret.username` | Container registry username                            | `""`      |
| `imagePullSecret.password` | Container registry password                            | `""`      |
| `serviceAccount.create`    | Create a service account                               | `true`    |
| `serviceAccount.name`      | Service account name                                   | `cs`      |
| `tls.autoGenerated`        | Generate automatically self-signed TLS certificates    | `true`    |
| `tls.verify`               | Verify connection to external cluster                  | `false`   |
| `tls.cert`                 | TLS certificate                                        | `""`      |
| `tls.certKey`              | TLS certificate key                                    | `""`      |
| `tls.certCA`               | TLS certificate CA                                     | `""`      |

### CS metrics

| Name              | Description       | Value  |
| ----------------- | ----------------- | ------ |
| `metrics.enabled` | Enable cs metrics | `true` |

### Auth-center component parameters

| Name                                       | Description                                                | Value |
| ------------------------------------------ | ---------------------------------------------------------- | ----- |
| `auth-center.nodeSelector`                 | Template to specify the labels of nodes for pod assignment | `{}`  |
| `auth-center.replicas`                     | Number of replicas for the auth-center component           | `2`   |
| `auth-center.administrator.existingSecret` | Name of the existing secret with administrator credentials | `""`  |
| `auth-center.administrator.username`       | Administrator name                                         | `""`  |
| `auth-center.administrator.password`       | Administrator password                                     | `""`  |

### Policy-enforcer component parameters

| Name                           | Description                                                | Value |
| ------------------------------ | ---------------------------------------------------------- | ----- |
| `policy-enforcer.nodeSelector` | Template to specify the labels of nodes for pod assignment | `{}`  |
| `policy-enforcer.replicas`     | Number of replicas for the policy-enforcer component       | `2`   |

### History-api component parameters

| Name                            | Description                                                | Value            |
| ------------------------------- | ---------------------------------------------------------- | ---------------- |
| `history-api.nodeSelector`      | Template to specify the labels of nodes for pod assignment | `{}`             |
| `history-api.replicas`          | Number of replicas for the history-api component           | `2`              |
| `history-api.retentionInterval` | Interval to retain history data for                        | `8760h`          |
| `history-api.rabbitmq.queue`    | RabbitMQ queue name                                        | `history_events` |

### Container-registry-integrator component parameters

| Name                           | Description                                                | Value |
| ------------------------------ | ---------------------------------------------------------- | ----- |
| `cluster-manager.nodeSelector` | Template to specify the labels of nodes for pod assignment | `{}`  |
| `cluster-manager.replicas`     | Number of replicas for the cluster-manager component       | `2`   |

### Notifier component parameters

| Name                    | Description                                                | Value |
| ----------------------- | ---------------------------------------------------------- | ----- |
| `notifier.nodeSelector` | Template to specify the labels of nodes for pod assignment | `{}`  |
| `notifier.replicas`     | Number of replicas for the notifier component              | `2`   |

### Reverse-proxy component parameters

| Name                                       | Description                                                       | Value   |
| ------------------------------------------ | ----------------------------------------------------------------- | ------- |
| `reverse-proxy.nodeSelector`               | Template to specify the labels of nodes for pod assignment        | `{}`    |
| `reverse-proxy.replicas`                   | Number of replicas for the reverse-proxy component                | `2`     |
| `reverse-proxy.ingress.enabled`            | Enable ingress for cs                                             | `false` |
| `reverse-proxy.ingress.class`              | Ingress class                                                     | `""`    |
| `reverse-proxy.ingress.hostname`           | Hostname of ingress                                               | `""`    |
| `reverse-proxy.ingress.tls.autoGenerated`  | Generate automatically self-signed TLS certificates               | `true`  |
| `reverse-proxy.ingress.tls.existingSecret` | Name of an existing secret that contains the certificates         | `""`    |
| `reverse-proxy.ingress.tls.cert`           | Certificate value. Requires `tls.autoGenerated` to be `false`     | `""`    |
| `reverse-proxy.ingress.tls.certKey`        | Certificate key value. Requires `tls.autoGenerated` to be `false` | `""`    |
| `reverse-proxy.ingress.tls.certCA`         | CA Certificate value. Requires `tls.autoGenerated` to be `false`  | `""`    |

### CS-manager component parameters

| Name                           | Description                                                | Value |
| ------------------------------ | ---------------------------------------------------------- | ----- |
| `cs-manager.nodeSelector`      | Template to specify the labels of nodes for pod assignment | `{}`  |
| `cs-manager.replicas`          | Number of replicas for the cs-manager component            | `1`   |
| `cs-manager.registrationToken` | Token for cluster registration                             | `""`  |

### Runtime-monitor component parameters

| Name                                         | Description                                                     | Value                                                                                                                                 |
| -------------------------------------------- | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `runtime-monitor.nodeSelector`               | Template to specify the labels of nodes for pod assignment      | `{}`                                                                                                                                  |
| `runtime-monitor.configUpdateInterval`       | Interval to update the stored configuration                     | `30s`                                                                                                                                 |
| `runtime-monitor.dnsPolicy`                  | Configuration of the DNS policy for runtime monitoring          | `ClusterFirstWithHostNet`                                                                                                             |
| `runtime-monitor.containerPorts.http`        | Port that HTTP server should be listening on                    | `9000`                                                                                                                                |
| `runtime-monitor.containerPorts.grpc`        | Port that GRPC server should be listening on                    | `8000`                                                                                                                                |
| `runtime-monitor.containerPorts.gops`        | Port that gops agent should be listening on                     | `7000`                                                                                                                                |
| `runtime-monitor.tetragon.enableProcessCred` | Enable visibility of capabilities in the exec and kprobe events | `true`                                                                                                                                |
| `runtime-monitor.tetragon.enableProcessNs`   | Enable visibility of namespaces in the exec and kprobe events   | `true`                                                                                                                                |
| `runtime-monitor.tetragon.exportAllowList`   | Allowlist for JSON export                                       | `{"pod_regex":["deathstar"],"event_set":["PROCESS_EXEC", "PROCESS_EXIT", "PROCESS_KPROBE", "PROCESS_UPROBE", "PROCESS_TRACEPOINT"]}
` |
| `runtime-monitor.tetragon.grpc.address`      | Set address of Tetragon grpc connection in host:port format     | `localhost:54321`                                                                                                                     |
| `runtime-monitor.tetragon.resources`         | Resource configuration for tetragon container                   | `{}`                                                                                                                                  |
| `runtime-monitor.rabbitmq.queue`             | RabbitMQ queue name                                             | `runtime_events`                                                                                                                      |
| `runtime-monitor.resources`                  | Resource configuration for runtime-monitor container            | `{}`                                                                                                                                  |

### Event-processor component parameters

| Name                                          | Description                                                | Value            |
| --------------------------------------------- | ---------------------------------------------------------- | ---------------- |
| `event-processor.nodeSelector`                | Template to specify the labels of nodes for pod assignment | `{}`             |
| `event-processor.replicas`                    | Number of replicas for the component                       | `2`              |
| `event-processor.configUpdateInterval`        | Interval to update the stored configuration                | `30s`            |
| `event-processor.rabbitmq.runtimeEventsQueue` | RabbitMQ runtime events queue name                         | `runtime_events` |
| `event-processor.rabbitmq.historyEventsQueue` | RabbitMQ history events queue name                         | `history_events` |
| `event-processor.resources`                   | Resource configuration for event-processor container       | `{}`             |

### Public-api component parameters

| Name                      | Description                                                | Value |
| ------------------------- | ---------------------------------------------------------- | ----- |
| `public-api.nodeSelector` | Template to specify the labels of nodes for pod assignment | `{}`  |
| `public-api.replicas`     | Number of replicas for the public-api component            | `2`   |

### Postgresql installation configuration

| Name                                        | Description                                                                                                                                                                                                                                                                                                                                          | Value               |
| ------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `postgresql.externalHost`                   | External host with PostgreSQL. Requires setting `postgresql.deploy` to `false`.                                                                                                                                                                                                                                                                      | `""`                |
| `postgresql.fullnameOverride`               | String to fully override common.names.fullname template                                                                                                                                                                                                                                                                                              | `postgresql`        |
| `postgresql.tls.autoGenerated`              | Generate automatically self-signed TLS certificates                                                                                                                                                                                                                                                                                                  | `true`              |
| `postgresql.tls.cert`                       | Certificate value. Requires `tls.autoGenerated` to be `false`                                                                                                                                                                                                                                                                                        | `""`                |
| `postgresql.tls.certKey`                    | Certificate key value. Requires `tls.autoGenerated` to be `false`                                                                                                                                                                                                                                                                                    | `""`                |
| `postgresql.tls.certCA`                     | CA Certificate value. Requires `tls.autoGenerated` to be `false`                                                                                                                                                                                                                                                                                     | `""`                |
| `postgresql.auth.existingSecret`            | Name of the existing secret with PostgreSQL credentials. The `auth.postgresPassword`, `auth.password`, and `auth.replicationPassword` values will be ignored and taken from this secret. The secret might also contain the `ldap-password` key if LDAP is enabled. If so, the `ldap.bind_password` value will be ignored and taken from this secret. | `postgresql`        |
| `postgresql.auth.username`                  | Name of custom user to be created                                                                                                                                                                                                                                                                                                                    | `cs`                |
| `postgresql.auth.password`                  | Password of custom user to be created. Ignored if `auth.existingSecret` is set.                                                                                                                                                                                                                                                                      | `""`                |
| `postgresql.auth.database`                  | Name of custom database to be created                                                                                                                                                                                                                                                                                                                | `cs`                |
| `postgresql.auth.existingSecretPasswordKey` | Name of the key in the existing secret with PostgreSQL credentials. Only used if `auth.existingSecret` is set.                                                                                                                                                                                                                                       | `POSTGRES_PASSWORD` |
| `postgresql.nodeSelector`                   | Labels of nodes for primary PostgreSQL pod assignment                                                                                                                                                                                                                                                                                                | `{}`                |
| `postgresql.resources`                      | Resource configuration for PostgreSQL container                                                                                                                                                                                                                                                                                                      | `{}`                |
| `postgresql.persistence.enabled`            | Enable data persistence for primary PostgreSQL using PVC                                                                                                                                                                                                                                                                                             | `true`              |
| `postgresql.persistence.storageClass`       | Persistent volume storage class for primary PostgreSQL                                                                                                                                                                                                                                                                                               | `""`                |
| `postgresql.persistence.size`               | Persistent volume size for PostgreSQL                                                                                                                                                                                                                                                                                                                | `1Gi`               |
| `postgresql.persistence.existingClaim`      | Name of an existing PVC                                                                                                                                                                                                                                                                                                                              | `""`                |
| `postgresql.persistence.selector`           | Template to specify an existing persistent volume                                                                                                                                                                                                                                                                                                    | `{}`                |
| `postgresql.metrics.enabled`                | Start a prometheus exporter                                                                                                                                                                                                                                                                                                                          | `true`              |
| `postgresql.metrics.externalHost`           | PostgreSQL metrics external host                                                                                                                                                                                                                                                                                                                     | `""`                |

### Redis installation configuration

| Name                                   | Description                                                                                | Value            |
| -------------------------------------- | ------------------------------------------------------------------------------------------ | ---------------- |
| `redis.externalHost`                   | External host with Redis. Requires setting `redis.deploy` to `false`.                      | `""`             |
| `redis.fullnameOverride`               | String to fully override common.names.fullname                                             | `redis`          |
| `redis.tls.autoGenerated`              | Generate automatically self-signed TLS certificates                                        | `true`           |
| `redis.tls.cert`                       | Certificate value. Requires `tls.autoGenerated` to be `false`                              | `""`             |
| `redis.tls.certKey`                    | Certificate key value. Requires `tls.autoGenerated` to be `false`                          | `""`             |
| `redis.tls.certCA`                     | CA Certificate value. Requires `tls.autoGenerated` to be `false`                           | `""`             |
| `redis.auth.existingSecret`            | Name of the existing secret with Redis credentials                                         | `redis`          |
| `redis.auth.username`                  | Redis username                                                                             | `cs`             |
| `redis.auth.password`                  | Redis password                                                                             | `""`             |
| `redis.auth.existingSecretPasswordKey` | Password key to retrieve from the existing secret                                          | `REDIS_PASSWORD` |
| `redis.replicaCount`                   | Number of Redis master instances to deploy (experimental, requires additional configuring) | `1`              |
| `redis.nodeSelector`                   | Labels of nodes for Redis master pod assignment                                            | `{}`             |
| `redis.resources`                      | Resource configuration for Redis container                                                 | `{}`             |
| `redis.persistence.enabled`            | Enable persistence for Redis master nodes using PVC                                        | `false`          |
| `redis.persistence.storageClass`       | Persistent volume storage class                                                            | `""`             |
| `redis.persistence.size`               | Persistent volume size                                                                     | `1Gi`            |
| `redis.persistence.existingClaim`      | Use an existing PVC created manually                                                       | `""`             |
| `redis.persistence.selector`           | Template to specify additional labels for PVC                                              | `{}`             |

### RabbitMQ installation configuration

| Name                                      | Description                                                                                          | Value             |
| ----------------------------------------- | ---------------------------------------------------------------------------------------------------- | ----------------- |
| `rabbitmq.externalHost`                   | External host with RabbitMQ. Requires setting `rabbitmq.deploy` to `false`.                          | `""`              |
| `rabbitmq.fullnameOverride`               | String to fully override rabbitmq.fullname template                                                  | `rabbitmq`        |
| `rabbitmq.auth.username`                  | RabbitMQ application username                                                                        | `cs`              |
| `rabbitmq.auth.password`                  | RabbitMQ application password                                                                        | `""`              |
| `rabbitmq.auth.existingSecret`            | Existing secret with RabbitMQ credentials (must contain value for the `rabbitmq-password` parameter) | `rabbitmq`        |
| `rabbitmq.auth.existingSecretPasswordKey` | Password key to be retrieved from existing secret                                                    | `RABBIT_PASSWORD` |
| `rabbitmq.nodeSelector`                   | Template to specify the labels of nodes for pod assignment                                           | `{}`              |
| `rabbitmq.resources`                      | Resource configuration for RabbitMQ container                                                        | `{}`              |
| `rabbitmq.persistence.enabled`            | Enable RabbitMQ data persistence using PVC                                                           | `true`            |
| `rabbitmq.persistence.storageClass`       | Persistent volume storage class for RabbitMQ                                                         | `""`              |
| `rabbitmq.persistence.size`               | Persistent volume size for RabbitMQ                                                                  | `1Gi`             |
| `rabbitmq.persistence.existingClaim`      | Name of an existing PVC                                                                              | `""`              |
| `rabbitmq.persistence.selector`           | Template to specify an existing persistent volume                                                    | `{}`              |
| `rabbitmq.metrics.enabled`                | Enable exposing RabbitMQ metrics to be gathered by Prometheus                                        | `true`            |
| `rabbitmq.metrics.externalHost`           | RabbitMQ metrics external host                                                                       | `""`              |

### Clickhouse installation configuration

| Name                                        | Description                                                                     | Value                 |
| ------------------------------------------- | ------------------------------------------------------------------------------- | --------------------- |
| `clickhouse.externalHost`                   | External host with ClickHouse. Requires setting `clickhouse.deploy` to `false`. | `""`                  |
| `clickhouse.fullnameOverride`               | String to fully override common.names.fullname                                  | `clickhouse`          |
| `clickhouse.nodeSelector`                   | Labels of nodes for ClickHouse pod assignment                                   | `{}`                  |
| `clickhouse.replicaCount`                   | Number of ClickHouse replicas to deploy per shard                               | `1`                   |
| `clickhouse.resources`                      | Resource configuration for Clickhouse container                                 | `{}`                  |
| `clickhouse.persistence.enabled`            | Enable persistence using PVC                                                    | `true`                |
| `clickhouse.persistence.storageClass`       | Persistent volume storage class                                                 | `""`                  |
| `clickhouse.persistence.size`               | Data volume size                                                                | `5Gi`                 |
| `clickhouse.persistence.existingClaim`      | Name of an existing PVC                                                         | `""`                  |
| `clickhouse.persistence.selector`           | Template to specify an existing persistent volume                               | `{}`                  |
| `clickhouse.tls.autoGenerated`              | Generate automatically self-signed TLS certificates                             | `true`                |
| `clickhouse.tls.cert`                       | Certificate value. Requires `tls.autoGenerated` to be `false`                   | `""`                  |
| `clickhouse.tls.certKey`                    | Certificate key value. Requires `tls.autoGenerated` to be `false`               | `""`                  |
| `clickhouse.tls.certCA`                     | CA Certificate value. Requires `tls.autoGenerated` to be `false`                | `""`                  |
| `clickhouse.auth.username`                  | ClickHouse administrator name                                                   | `cs`                  |
| `clickhouse.auth.password`                  | ClickHouse administartor password                                               | `""`                  |
| `clickhouse.auth.existingSecret`            | Name of the secret with the administrator password                              | `clickhouse`          |
| `clickhouse.auth.existingSecretPasswordKey` | Name of the key stored in the existing secret                                   | `CLICKHOUSE_PASSWORD` |
| `clickhouse.auth.database`                  | Name of the ClickHouse database                                                 | `cs`                  |
| `clickhouse.metrics.enabled`                | Enable the export of Prometheus metrics                                         | `true`                |
| `clickhouse.metrics.externalHost`           | ClickHouse metrics external host                                                | `""`                  |
