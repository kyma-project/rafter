{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "rafterUploadService.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rafterUploadService.fullname" -}}
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

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "rafterUploadService.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service
*/}}
{{- define "rafterUploadService.serviceName" -}}
{{- if .Values.service.name -}}
{{- include .Values.service.name . | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- include "rafterUploadService.fullname" . -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account
*/}}
{{- define "rafterUploadService.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "rafterUploadService.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac cluster role
*/}}
{{- define "rafterUploadService.rbacClusterRoleName" -}}
{{- if .Values.rbac.clusterScope.create -}}
    {{ default (include "rafterUploadService.fullname" .) .Values.rbac.clusterScope.role.name }}
{{- else -}}
    {{ default "default" .Values.rbac.clusterScope.role.name  }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac cluster role binding
*/}}
{{- define "rafterUploadService.rbacClusterRoleBindingName" -}}
{{- if .Values.rbac.clusterScope.create -}}
    {{ default (include "rafterUploadService.fullname" .) .Values.rbac.clusterScope.roleBinding.name }}
{{- else -}}
    {{ default "default" .Values.rbac.clusterScope.roleBinding.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service monitor
*/}}
{{- define "rafterUploadService.serviceMonitorName" -}}
{{- if and .Values.serviceMonitor.create }}
    {{ default (include "rafterUploadService.fullname" .) .Values.serviceMonitor.name }}
{{- else -}}
    {{ default "default" .Values.serviceMonitor.name }}
{{- end -}}
{{- end -}}

{{/*
Renders a value that contains template.
Usage:
{{ include "rafterUploadService.tplValue" ( dict "value" .Values.path.to.the.Value "context" $ ) }}
*/}}
{{- define "rafterUploadService.tplValue" -}}
    {{- if typeIs "string" .value }}
        {{- tpl .value .context }}
    {{- else }}
        {{- tpl (.value | toYaml) .context }}
    {{- end }}
{{- end -}}

{{/*
Renders a proper env in container
Usage:
{{ include "rafterUploadService.createEnv" ( dict "name" "APP_FOO_BAR" "value" .Values.path.to.the.Value "context" $ ) }}
*/}}
{{- define "rafterUploadService.createEnv" -}}
{{- if and .name .value }}
{{- printf "- name: %s" .name -}}
{{- include "rafterUploadService.tplValue" ( dict "value" .value "context" .context ) | nindent 2 }}
{{- end }}
{{- end -}}
