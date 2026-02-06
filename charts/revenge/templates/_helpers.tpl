{{- define "revenge.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

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

{{- define "revenge.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{ include "revenge.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "revenge.selectorLabels" -}}
app.kubernetes.io/name: {{ include "revenge.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "revenge.databaseHost" -}}
{{- if .Values.postgresql.enabled }}
{{- printf "%s-postgresql" (include "revenge.fullname" .) }}
{{- else }}
{{- .Values.revenge.database.host }}
{{- end }}
{{- end }}

{{- define "revenge.cacheHost" -}}
{{- if .Values.dragonfly.enabled }}
{{- printf "%s-dragonfly" (include "revenge.fullname" .) }}
{{- else }}
{{- .Values.revenge.cache.host }}
{{- end }}
{{- end }}

{{- define "revenge.searchHost" -}}
{{- if .Values.typesense.enabled }}
{{- printf "%s-typesense" (include "revenge.fullname" .) }}
{{- else }}
{{- .Values.revenge.search.host }}
{{- end }}
{{- end }}

{{- define "revenge.databaseURL" -}}
postgres://{{ .Values.revenge.database.user }}:$(REVENGE_DATABASE_PASSWORD)@{{ include "revenge.databaseHost" . }}:{{ .Values.revenge.database.port }}/{{ .Values.revenge.database.name }}?sslmode=disable
{{- end }}
