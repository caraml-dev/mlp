{{/* vim: set filetype=mustache: */}}
{{- define "mlp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "mlp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version -}}
{{- end -}}

{{- define "mlp.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "postgres.host" -}}
{{ if .Values.postgresql.enabled }}
    {{- printf "%s-postgresql.%s.svc.cluster.local" .Release.Name .Release.Namespace -}}
{{- else -}}
    {{- .Values.externalPostgresql.address -}}
{{- end -}}
{{- end -}}

{{- define "postgres.username" -}}
{{- if .Values.postgresql.enabled -}}
{{- .Values.postgresql.postgresqlUsername -}}
{{- else -}}
{{- .Values.externalPostgresql.username -}}
{{- end -}}
{{- end -}}

{{- define "postgres.database" -}}
{{- if .Values.postgresql.enabled -}}
{{- .Values.postgresql.postgresqlDatabase -}}
{{- else -}}
{{- .Values.externalPostgresql.database -}}
{{- end -}}
{{- end -}}
