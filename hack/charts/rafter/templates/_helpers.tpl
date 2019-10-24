{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "rafter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rafter.fullname" -}}
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
{{- define "rafter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account
*/}}
{{- define "rafter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "rafter.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac role
*/}}
{{- define "rafter.rbacRoleName" -}}
{{- if .Values.rbac.namespaced.create -}}
    {{ default (include "rafter.fullname" .) .Values.rbac.namespaced.role.name }}
{{- else -}}
    {{ default "default" .Values.rbac.namespaced.role.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac role binding
*/}}
{{- define "rafter.rbacRoleBindingName" -}}
{{- if .Values.rbac.namespaced.create -}}
    {{ default (include "rafter.fullname" .) .Values.rbac.namespaced.roleBinding.name }}
{{- else -}}
    {{ default "default" .Values.rbac.namespaced.roleBinding.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac cluster role
*/}}
{{- define "rafter.rbacClusterRoleName" -}}
{{- if .Values.rbac.clusterScope.create -}}
    {{ default (include "rafter.fullname" .) .Values.rbac.clusterScope.role.name }}
{{- else -}}
    {{ default "default" .Values.rbac.clusterScope.role.name  }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the rbac cluster role binding
*/}}
{{- define "rafter.rbacClusterRoleBindingName" -}}
{{- if .Values.rbac.clusterScope.create -}}
    {{ default (include "rafter.fullname" .) .Values.rbac.clusterScope.roleBinding.name }}
{{- else -}}
    {{ default "default" .Values.rbac.clusterScope.roleBinding.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the config map with webhooks
*/}}
{{- define "rafter.webhooksConfigMapName" -}}
{{- if and .Values.webhooks.enabled -}}
    {{ default (printf "%s-%s" (include "rafter.fullname" .) "webhooks") .Values.webhooks.configMap.name }}
{{- else -}}
    {{ default "default" .Values.webhooks.configMap.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the metrics service
*/}}
{{- define "rafter.metricsServiceName" -}}
{{- if .Values.metrics.enabled -}}
    {{ default (include "rafter.fullname" .) .Values.metrics.service.name }}
{{- else -}}
    {{ default "default" .Values.metrics.service.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service monitor
*/}}
{{- define "rafter.serviceMonitorName" -}}
{{- if and .Values.metrics.enabled .Values.metrics.serviceMonitor.enabled }}
    {{ default (include "rafter.fullname" .) .Values.metrics.serviceMonitor.name }}
{{- else -}}
    {{ default "default" .Values.metrics.serviceMonitor.name }}
{{- end -}}
{{- end -}}
