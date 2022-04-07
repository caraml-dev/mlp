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
{{- printf "%s-postgresql.%s.svc.cluster.local" .Release.Name .Release.Namespace -}}
{{- end -}}

{{- define "postgres.username" -}}
{{- if .Values.externalPostgres.enabled -}}
{{- .Values.externalPostgres.username -}}
{{- else -}}
{{- .Values.postgresql.postgresqlUsername -}}
{{- end -}}
{{- end -}}

{{- define "postgres.database" -}}
{{- if .Values.externalPostgres.enabled -}}
{{- .Values.externalPostgres.database -}}
{{- else -}}
{{- .Values.postgresql.postgresqlDatabase -}}
{{- end -}}
{{- end -}}
