# Revenge Helm Chart

Deploys Revenge media server to Kubernetes with support for clustered deployments.

## Prerequisites

- Kubernetes 1.23+
- Helm 3.8+
- PostgreSQL 18+ (can use included subchart or external)
- Dragonfly cache (Redis-compatible)
- Typesense 0.25+ for search
- NFS server or compatible storage for media files (optional but recommended for clusters)

## Media Storage Configuration

Revenge requires shared media storage when running multiple replicas. Three configuration options are supported:

### Option 1: Direct NFS Mount (Recommended)

Mount NFS directly in pods without creating PV/PVC:

```yaml
media:
  persistence:
    enabled: true
    nfs:
      server: "nas.example.com"
      path: "/volume1/media"
      readOnly: true
```

**Pros**: Simple, no PV/PVC needed
**Cons**: NFS credentials in values file

### Option 2: Existing PersistentVolumeClaim

Use a pre-created PVC:

```yaml
media:
  persistence:
    enabled: true
    existingClaim: "my-media-pvc"
```

**Pros**: Separates storage provisioning from app deployment
**Cons**: Requires manual PVC creation

### Option 3: Automatic PVC with StorageClass

Let Helm create a PVC using a StorageClass:

```yaml
media:
  persistence:
    enabled: true
    storageClass: "nfs-client"  # or your StorageClass
    accessMode: ReadOnlyMany
    size: 1Ti
```

**Pros**: Declarative, GitOps-friendly
**Cons**: Requires NFS provisioner or compatible StorageClass

## Installation

### Basic installation with external database

```bash
helm install revenge ./charts/revenge \
  --set revenge.database.host=postgres.example.com \
  --set revenge.database.password=secretpassword \
  --set revenge.cache.host=dragonfly.example.com \
  --set media.persistence.enabled=true \
  --set media.persistence.nfs.server=nas.example.com \
  --set media.persistence.nfs.path=/volume1/media
```

### With included PostgreSQL subchart

```bash
helm install revenge ./charts/revenge \
  --set postgresql.enabled=true \
  --set media.persistence.enabled=true \
  --set media.persistence.nfs.server=nas.example.com \
  --set media.persistence.nfs.path=/volume1/media
```

### Using values file

Create `my-values.yaml`:

```yaml
replicaCount: 3

revenge:
  database:
    host: postgres.prod.svc.cluster.local
    password: changeme
  cache:
    host: dragonfly.prod.svc.cluster.local

media:
  persistence:
    enabled: true
    nfs:
      server: nas.example.com
      path: /volume1/media
      readOnly: true
  moviePaths:
    - /media/movies
  tvPaths:
    - /media/tv

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: revenge.example.com
      paths:
        - path: /
          pathType: Prefix
```

Install:

```bash
helm install revenge ./charts/revenge -f my-values.yaml
```

## Media Path Configuration

Configure which paths inside the container to use for each media type:

```yaml
media:
  moviePaths:
    - /media/movies
    - /media/films  # Multiple paths supported
  tvPaths:
    - /media/tv
  musicPaths:
    - /media/music
```

These paths should exist within the mounted `/media` volume.

## Storage Class Setup

If using automatic PVC with a StorageClass, ensure you have an NFS provisioner installed:

### Using nfs-subdir-external-provisioner

```bash
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
helm install nfs-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
  --set nfs.server=nas.example.com \
  --set nfs.path=/volume1/k8s
```

Then use `storageClass: nfs-client` in your values.

## Examples

See `examples/` directory for:
- `media-pv-nfs.yaml` - Example PersistentVolume for NFS
- Additional configuration examples

## Upgrading

```bash
helm upgrade revenge ./charts/revenge -f my-values.yaml
```

## Uninstalling

```bash
helm uninstall revenge
```

**Note**: PersistentVolumeClaims are not deleted automatically. Delete manually if needed:

```bash
kubectl delete pvc revenge-media
```

## Configuration

See `values.yaml` for all available configuration options.

### Key Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `2` |
| `media.persistence.enabled` | Enable media volume | `false` |
| `media.persistence.storageClass` | StorageClass for PVC | `nfs-client` |
| `media.persistence.size` | Size of media volume | `1Ti` |
| `media.persistence.existingClaim` | Use existing PVC | `""` |
| `media.persistence.nfs.server` | NFS server address | `""` |
| `media.persistence.nfs.path` | NFS export path | `""` |
| `autoscaling.enabled` | Enable HPA | `true` |
| `autoscaling.minReplicas` | Minimum replicas | `2` |
| `autoscaling.maxReplicas` | Maximum replicas | `10` |

## Troubleshooting

### Media not accessible

Check if media volume is mounted:

```bash
kubectl exec -it deployment/revenge -- ls -la /media
```

### NFS mount fails

Check NFS server connectivity:

```bash
kubectl exec -it deployment/revenge -- ping nas.example.com
```

Check NFS exports on server:

```bash
showmount -e nas.example.com
```

### Multiple replicas can't access media

Ensure `ReadOnlyMany` access mode is used and NFS server supports it.

## Support

For issues and questions, see the main Revenge repository.
