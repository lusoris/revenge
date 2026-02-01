# Deployment

This directory contains deployment configurations for running Revenge in various environments.

> **Note**: These files are auto-generated from SOURCE_OF_TRUTH.md.
> Do not edit directly. Run `python scripts/deploy-pipeline/01-generate.py` to regenerate.

## Docker Compose (Development / Single Node)

```bash
docker compose -f deploy/docker-compose.yml up -d
```

## Docker Swarm (Production)

```bash
# Create secrets
echo "your-password" | docker secret create db_password -

# Deploy
docker stack deploy -c deploy/docker-swarm-stack.yml revenge
```

## Kubernetes (Helm)

```bash
helm install revenge charts/revenge --namespace revenge --create-namespace
```

For OCI chart:
```bash
helm install revenge oci://ghcr.io/lusoris/charts/revenge
```

## Environment Variables

See `docs/dev/design/00_SOURCE_OF_TRUTH.md` for complete configuration reference.
