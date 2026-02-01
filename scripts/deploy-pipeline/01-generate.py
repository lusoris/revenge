#!/usr/bin/env python3
"""Deploy File Generator - Generate deployment configs from SOURCE_OF_TRUTH.

This script generates:
1. Helm chart (charts/revenge/)
2. Docker Compose file (deploy/docker-compose.yml)
3. Docker Swarm stack (deploy/docker-swarm-stack.yml)
4. Deploy README (deploy/README.md)

All configs are generated using templates + data from SOURCE_OF_TRUTH.md

Author: Automation System
Created: 2026-02-01
"""

import re
import sys
from pathlib import Path
from typing import Any

import yaml


# Add parent to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent))
from automation.sot_parser import SOTParser


class DeployGenerator:
    """Generate deployment configs from SOURCE_OF_TRUTH."""

    def __init__(self, repo_root: Path):
        """Initialize generator."""
        self.repo_root = repo_root
        self.sot_path = repo_root / "docs" / "dev" / "design" / "00_SOURCE_OF_TRUTH.md"
        self.templates_dir = repo_root / "templates" / "deploy"
        self.charts_dir = repo_root / "charts" / "revenge"
        self.deploy_dir = repo_root / "deploy"

        # Parse SOT
        self.sot = self._parse_sot()

    def _parse_sot(self) -> dict[str, Any]:
        """Parse SOURCE_OF_TRUTH for deployment data."""
        if not self.sot_path.exists():
            print(f"❌ Error: {self.sot_path} not found")
            sys.exit(1)

        parser = SOTParser(self.sot_path)
        data = parser.parse()

        # Also extract deployment-specific sections
        content = self.sot_path.read_text()
        data["deploy"] = self._parse_deploy_sections(content)

        return data

    def _parse_deploy_sections(self, content: str) -> dict[str, Any]:
        """Parse deployment-specific sections from SOT."""
        deploy = {}

        # Extract K8s resources section
        k8s_match = re.search(
            r"### Kubernetes Resources\n\n```yaml\n(.*?)```",
            content,
            re.DOTALL,
        )
        if k8s_match:
            try:
                deploy["k8s_example"] = yaml.safe_load(k8s_match.group(1))
            except yaml.YAMLError:
                deploy["k8s_example"] = {}

        # Extract Docker Swarm section
        swarm_match = re.search(
            r"### Docker Swarm Stack\n\n```yaml\n(.*?)```",
            content,
            re.DOTALL,
        )
        if swarm_match:
            try:
                deploy["swarm_example"] = yaml.safe_load(swarm_match.group(1))
            except yaml.YAMLError:
                deploy["swarm_example"] = {}

        # Extract Helm chart structure
        helm_match = re.search(
            r"### Helm Chart Structure\n\n```\n(.*?)```",
            content,
            re.DOTALL,
        )
        if helm_match:
            deploy["helm_structure"] = helm_match.group(1).strip()

        # Extract scaling strategy
        scaling_match = re.search(
            r"### Scaling Strategy\n\n\|.*?\n\|.*?\n(.*?)(?=\n###|\n##|\Z)",
            content,
            re.DOTALL,
        )
        if scaling_match:
            deploy["scaling"] = self._parse_scaling_table(scaling_match.group(0))

        # Extract health probes
        probes_match = re.search(
            r"### Health Endpoints\n\n\|.*?\n\|.*?\n(.*?)(?=\n###|\n##|\Z)",
            content,
            re.DOTALL,
        )
        if probes_match:
            deploy["probes"] = self._parse_probes_table(probes_match.group(0))

        # Extract config keys for env mapping
        config_match = re.search(
            r"### Server Configuration\n\n\|.*?\n\|.*?\n(.*?)(?=\n###|\n##|\Z)",
            content,
            re.DOTALL,
        )
        if config_match:
            deploy["config_keys"] = self._parse_config_table(config_match.group(0))

        return deploy

    def _parse_scaling_table(self, table: str) -> list[dict]:
        """Parse scaling strategy table."""
        scaling = []
        lines = [l.strip() for l in table.split("\n") if "|" in l]
        for line in lines[2:]:  # Skip header and separator
            cells = [c.strip() for c in line.split("|")[1:-1]]
            if len(cells) >= 4:
                scaling.append(
                    {
                        "component": cells[0],
                        "min": cells[1],
                        "max": cells[2],
                        "trigger": cells[3],
                    }
                )
        return scaling

    def _parse_probes_table(self, table: str) -> list[dict]:
        """Parse health probes table."""
        probes = []
        lines = [l.strip() for l in table.split("\n") if "|" in l]
        for line in lines[2:]:
            cells = [c.strip() for c in line.split("|")[1:-1]]
            if len(cells) >= 3:
                probes.append(
                    {
                        "endpoint": cells[0].replace("`", ""),
                        "purpose": cells[1],
                        "status_code": cells[2],
                    }
                )
        return probes

    def _parse_config_table(self, table: str) -> list[dict]:
        """Parse configuration keys table."""
        config = []
        lines = [l.strip() for l in table.split("\n") if "|" in l]
        for line in lines[2:]:
            cells = [c.strip() for c in line.split("|")[1:-1]]
            if len(cells) >= 4:
                config.append(
                    {
                        "key": cells[0].replace("`", ""),
                        "env": cells[1].replace("`", ""),
                        "type": cells[2],
                        "default": cells[3],
                    }
                )
        return config

    def generate_all(self) -> None:
        """Generate all deployment configs."""
        print("\n" + "=" * 70)
        print("DEPLOYMENT FILE GENERATOR")
        print("=" * 70 + "\n")

        self._generate_helm_chart()
        self._generate_docker_compose()
        self._generate_docker_swarm()
        self._generate_readme()

        print("\n" + "=" * 70)
        print("✅ Deployment files generated!")
        print("=" * 70 + "\n")

    def _generate_helm_chart(self) -> None:
        """Generate Helm chart from templates."""
        print("📦 Generating Helm chart...")

        # Create directories
        templates_dir = self.charts_dir / "templates"
        templates_dir.mkdir(parents=True, exist_ok=True)

        # Get metadata from SOT (versions used in Chart.yaml appVersion etc.)
        _ = self.sot.get("metadata", {})

        # Generate Chart.yaml
        chart = {
            "apiVersion": "v2",
            "name": "revenge",
            "description": "A modern, self-hosted media server with Go backend, SvelteKit frontend, and PostgreSQL database",
            "type": "application",
            "version": "0.1.0",
            "appVersion": "0.1.0",
            "home": "https://github.com/lusoris/revenge",
            "sources": ["https://github.com/lusoris/revenge"],
            "maintainers": [
                {"name": "Revenge Team", "url": "https://github.com/lusoris/revenge"}
            ],
            "keywords": ["media-server", "self-hosted", "go", "svelte"],
            # Note: Dependencies like postgresql can be added by users via:
            # helm dependency add bitnami/postgresql --repository oci://registry-1.docker.io/bitnamicharts
        }

        self._write_yaml(self.charts_dir / "Chart.yaml", chart)

        # Generate values.yaml with SOT data
        self._generate_helm_values()

        # Generate template files
        self._generate_helm_templates()

        # Generate .helmignore
        helmignore = """# Patterns to ignore when building packages
.DS_Store
*.swp
*.bak
*.tmp
*~
.git
.gitignore
.vscode
"""
        (self.charts_dir / ".helmignore").write_text(helmignore)

        print(f"   ✓ Helm chart generated at {self.charts_dir}")

    def _generate_helm_values(self) -> None:
        """Generate Helm values.yaml from SOT."""
        deploy = self.sot.get("deploy", {})

        # Get probes from SOT
        probes = deploy.get("probes", [])
        liveness_path = "/health/live"
        readiness_path = "/health/ready"
        for probe in probes:
            if "Liveness" in probe.get("purpose", ""):
                liveness_path = probe.get("endpoint", liveness_path)
            elif "Readiness" in probe.get("purpose", ""):
                readiness_path = probe.get("endpoint", readiness_path)

        # Get scaling from SOT
        scaling = deploy.get("scaling", [])
        min_replicas = 2
        max_replicas = 10
        cpu_target = 70
        for scale in scaling:
            if "API" in scale.get("component", ""):
                min_replicas = int(scale.get("min", 2))
                max_replicas = int(scale.get("max", 10))
                trigger = scale.get("trigger", "")
                cpu_match = re.search(r"CPU\s*>\s*(\d+)", trigger)
                if cpu_match:
                    cpu_target = int(cpu_match.group(1))

        values = {
            "replicaCount": min_replicas,
            "image": {
                "repository": "ghcr.io/lusoris/revenge",
                "pullPolicy": "IfNotPresent",
                "tag": "",
            },
            "serviceAccount": {
                "create": True,
                "automount": True,
                "annotations": {},
                "name": "",
            },
            "service": {
                "type": "ClusterIP",
                "port": 8080,
            },
            "ingress": {
                "enabled": False,
                "className": "",
                "annotations": {},
                "hosts": [
                    {
                        "host": "revenge.local",
                        "paths": [{"path": "/", "pathType": "Prefix"}],
                    }
                ],
                "tls": [],
            },
            "resources": {
                "limits": {"memory": "1Gi"},
                "requests": {"cpu": "100m", "memory": "256Mi"},
            },
            "livenessProbe": {
                "httpGet": {"path": liveness_path, "port": "http"},
                "initialDelaySeconds": 10,
                "periodSeconds": 30,
                "timeoutSeconds": 10,
            },
            "readinessProbe": {
                "httpGet": {"path": readiness_path, "port": "http"},
                "initialDelaySeconds": 5,
                "periodSeconds": 10,
                "timeoutSeconds": 5,
            },
            "autoscaling": {
                "enabled": True,
                "minReplicas": min_replicas,
                "maxReplicas": max_replicas,
                "targetCPUUtilizationPercentage": cpu_target,
            },
            "revenge": {
                "server": {"port": 8080, "host": "0.0.0.0"},
                "database": {
                    "host": "",
                    "port": 5432,
                    "name": "revenge",
                    "user": "revenge",
                },
                "cache": {"host": "", "port": 6379},
            },
            "postgresql": {"enabled": True},
        }

        self._write_yaml(self.charts_dir / "values.yaml", values)

    def _generate_helm_templates(self) -> None:
        """Generate Helm template files."""
        templates_dir = self.charts_dir / "templates"

        # _helpers.tpl
        helpers = """{{/*
Expand the name of the chart.
*/}}
{{- define "revenge.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "revenge.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "revenge.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{ include "revenge.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "revenge.selectorLabels" -}}
app.kubernetes.io/name: {{ include "revenge.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "revenge.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "revenge.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database host
*/}}
{{- define "revenge.databaseHost" -}}
{{- if .Values.postgresql.enabled }}
{{- printf "%s-postgresql" (include "revenge.fullname" .) }}
{{- else }}
{{- .Values.revenge.database.host }}
{{- end }}
{{- end }}
"""
        (templates_dir / "_helpers.tpl").write_text(helpers)

        # deployment.yaml (simplified)
        deployment = """apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "revenge.fullname" . }}
  labels:
    {{- include "revenge.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "revenge.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "revenge.labels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: {{ .Values.revenge.server.port }}
          env:
            - name: REVENGE_SERVER_PORT
              value: {{ .Values.revenge.server.port | quote }}
            - name: REVENGE_DATABASE_HOST
              value: {{ include "revenge.databaseHost" . | quote }}
          livenessProbe:
            {{- toYaml .Values.livenessProbe | nindent 12 }}
          readinessProbe:
            {{- toYaml .Values.readinessProbe | nindent 12 }}
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
"""
        (templates_dir / "deployment.yaml").write_text(deployment)

        # service.yaml
        service = """apiVersion: v1
kind: Service
metadata:
  name: {{ include "revenge.fullname" . }}
  labels:
    {{- include "revenge.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      name: http
  selector:
    {{- include "revenge.selectorLabels" . | nindent 4 }}
"""
        (templates_dir / "service.yaml").write_text(service)

        # hpa.yaml
        hpa = """{{- if .Values.autoscaling.enabled }}
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "revenge.fullname" . }}
  labels:
    {{- include "revenge.labels" . | nindent 4 }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "revenge.fullname" . }}
  minReplicas: {{ .Values.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.autoscaling.maxReplicas }}
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
{{- end }}
"""
        (templates_dir / "hpa.yaml").write_text(hpa)

    def _generate_docker_compose(self) -> None:
        """Generate Docker Compose file from SOT."""
        print("🐳 Generating Docker Compose...")

        self.deploy_dir.mkdir(parents=True, exist_ok=True)

        meta = self.sot.get("metadata", {})

        # Get PostgreSQL version from metadata (e.g., "18.1 (ONLY - no SQLite support)")
        pg_version_raw = meta.get("postgresql_version", "18")
        pg_version = re.match(r"(\d+)", pg_version_raw)
        pg_version = pg_version.group(1) if pg_version else "18"

        compose = {
            "version": "3.8",
            "services": {
                "revenge": {
                    "image": "ghcr.io/lusoris/revenge:${REVENGE_VERSION:-latest}",
                    "build": {"context": "..", "dockerfile": "Dockerfile"},
                    "ports": ["${REVENGE_PORT:-8080}:8080"],
                    "environment": {
                        "REVENGE_SERVER_PORT": "8080",
                        "REVENGE_SERVER_HOST": "0.0.0.0",
                        "REVENGE_DATABASE_HOST": "db",
                        "REVENGE_DATABASE_PORT": "5432",
                        "REVENGE_DATABASE_NAME": "revenge",
                        "REVENGE_DATABASE_USER": "revenge",
                        "REVENGE_DATABASE_PASSWORD": "${DB_PASSWORD:-revenge_dev_password}",
                        "REVENGE_CACHE_HOST": "dragonfly",
                        "REVENGE_CACHE_PORT": "6379",
                    },
                    "volumes": ["${MEDIA_PATH:-./media}:/media:ro"],
                    "networks": ["revenge-net"],
                    "depends_on": {
                        "db": {"condition": "service_healthy"},
                        "dragonfly": {"condition": "service_healthy"},
                    },
                    "healthcheck": {
                        "test": [
                            "CMD",
                            "wget",
                            "--no-verbose",
                            "--tries=1",
                            "--spider",
                            "http://localhost:8080/health/ready",
                        ],
                        "interval": "30s",
                        "timeout": "10s",
                        "retries": 3,
                    },
                    "restart": "unless-stopped",
                },
                "db": {
                    "image": f"postgres:{pg_version}-alpine",
                    "environment": {
                        "POSTGRES_DB": "revenge",
                        "POSTGRES_USER": "revenge",
                        "POSTGRES_PASSWORD": "${DB_PASSWORD:-revenge_dev_password}",
                    },
                    "volumes": ["pgdata:/var/lib/postgresql/data"],
                    "networks": ["revenge-net"],
                    "healthcheck": {
                        "test": ["CMD-SHELL", "pg_isready -U revenge -d revenge"],
                        "interval": "10s",
                        "timeout": "5s",
                        "retries": 5,
                    },
                    "restart": "unless-stopped",
                },
                "dragonfly": {
                    "image": "docker.dragonflydb.io/dragonflydb/dragonfly:latest",
                    "command": "dragonfly --logtostderr",
                    "volumes": ["dfdata:/data"],
                    "networks": ["revenge-net"],
                    "healthcheck": {
                        "test": ["CMD", "redis-cli", "ping"],
                        "interval": "10s",
                        "timeout": "5s",
                        "retries": 5,
                    },
                    "restart": "unless-stopped",
                },
            },
            "volumes": {"pgdata": None, "dfdata": None},
            "networks": {"revenge-net": {"driver": "bridge"}},
        }

        self._write_yaml(self.deploy_dir / "docker-compose.yml", compose)
        print(
            f"   ✓ Docker Compose generated at {self.deploy_dir / 'docker-compose.yml'}"
        )

    def _generate_docker_swarm(self) -> None:
        """Generate Docker Swarm stack from SOT."""
        print("🐝 Generating Docker Swarm stack...")

        meta = self.sot.get("metadata", {})
        deploy = self.sot.get("deploy", {})

        # Get PostgreSQL version from metadata (e.g., "18.1 (ONLY - no SQLite support)")
        pg_version_raw = meta.get("postgresql_version", "18")
        pg_version = re.match(r"(\d+)", pg_version_raw)
        pg_version = pg_version.group(1) if pg_version else "18"

        # Get scaling
        replicas = 2
        scaling = deploy.get("scaling", [])
        for scale in scaling:
            if "API" in scale.get("component", ""):
                replicas = int(scale.get("min", 2))

        swarm = {
            "version": "3.8",
            "services": {
                "revenge": {
                    "image": "ghcr.io/lusoris/revenge:${REVENGE_VERSION:-latest}",
                    "deploy": {
                        "replicas": replicas,
                        "update_config": {
                            "parallelism": 1,
                            "delay": "10s",
                            "failure_action": "rollback",
                        },
                        "restart_policy": {
                            "condition": "on-failure",
                            "max_attempts": 3,
                        },
                        "resources": {
                            "limits": {"memory": "1G"},
                            "reservations": {"memory": "256M"},
                        },
                    },
                    "environment": {
                        "REVENGE_SERVER_PORT": "8080",
                        "REVENGE_DATABASE_HOST": "db",
                        "REVENGE_DATABASE_PORT": "5432",
                        "REVENGE_DATABASE_NAME": "revenge",
                        "REVENGE_DATABASE_USER": "revenge",
                        "REVENGE_DATABASE_PASSWORD_FILE": "/run/secrets/db_password",
                        "REVENGE_CACHE_HOST": "dragonfly",
                        "REVENGE_CACHE_PORT": "6379",
                    },
                    "secrets": ["db_password"],
                    "volumes": ["media:/media:ro"],
                    "networks": ["revenge-net"],
                    "healthcheck": {
                        "test": [
                            "CMD",
                            "wget",
                            "--no-verbose",
                            "--tries=1",
                            "--spider",
                            "http://localhost:8080/health/ready",
                        ],
                        "interval": "30s",
                        "timeout": "10s",
                        "retries": 3,
                    },
                },
                "db": {
                    "image": f"postgres:{pg_version}-alpine",
                    "deploy": {
                        "replicas": 1,
                        "placement": {"constraints": ["node.role == manager"]},
                    },
                    "environment": {
                        "POSTGRES_DB": "revenge",
                        "POSTGRES_USER": "revenge",
                        "POSTGRES_PASSWORD_FILE": "/run/secrets/db_password",
                    },
                    "secrets": ["db_password"],
                    "volumes": ["pgdata:/var/lib/postgresql/data"],
                    "networks": ["revenge-net"],
                },
                "dragonfly": {
                    "image": "docker.dragonflydb.io/dragonflydb/dragonfly:latest",
                    "deploy": {
                        "replicas": 1,
                        "placement": {"constraints": ["node.role == manager"]},
                    },
                    "command": "dragonfly --logtostderr",
                    "volumes": ["dfdata:/data"],
                    "networks": ["revenge-net"],
                },
            },
            "secrets": {"db_password": {"external": True}},
            "volumes": {
                "media": {
                    "driver": "local",
                    "driver_opts": {
                        "type": "nfs",
                        "o": "addr=${NFS_SERVER:-nas.local},rw",
                        "device": ":${NFS_PATH:-/volume1/media}",
                    },
                },
                "pgdata": None,
                "dfdata": None,
            },
            "networks": {"revenge-net": {"driver": "overlay"}},
        }

        self._write_yaml(self.deploy_dir / "docker-swarm-stack.yml", swarm)
        print(
            f"   ✓ Docker Swarm stack generated at {self.deploy_dir / 'docker-swarm-stack.yml'}"
        )

    def _generate_readme(self) -> None:
        """Generate deploy README from SOT."""
        print("📝 Generating deploy README...")

        readme = """# Deployment

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
"""

        (self.deploy_dir / "README.md").write_text(readme)
        print(f"   ✓ README generated at {self.deploy_dir / 'README.md'}")

    def _write_yaml(self, path: Path, data: dict) -> None:
        """Write YAML file with comment header."""
        path.parent.mkdir(parents=True, exist_ok=True)

        header = f"""# Auto-generated from SOURCE_OF_TRUTH.md
# Do not edit directly. Run: python scripts/deploy-pipeline/01-generate.py
# Generated: {__import__("datetime").datetime.now().strftime("%Y-%m-%d")}

"""
        with open(path, "w") as f:
            f.write(header)
            yaml.dump(
                data, f, default_flow_style=False, sort_keys=False, allow_unicode=True
            )


def main():
    """Main entry point."""
    repo_root = Path(__file__).parent.parent.parent

    generator = DeployGenerator(repo_root)
    generator.generate_all()


if __name__ == "__main__":
    main()
