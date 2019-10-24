{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "rafterAsyncAPIService.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rafterAsyncAPIService.fullname" -}}
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
{{- define "rafterAsyncAPIService.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service
*/}}
{{- define "rafterAsyncAPIService.serviceName" -}}
{{- if .Values.service.name -}}
{{- printf "%s-%s" (include "rafterAsyncAPIService.fullname" .) .Values.service.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- include "rafterAsyncAPIService.fullname" . | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service monitor
*/}}
{{- define "rafterAsyncAPIService.serviceMonitorName" -}}
{{- if and .Values.metrics.enabled }}
    {{ default (include "rafterAsyncAPIService.fullname" .) .Values.metrics.serviceMonitor.name }}
{{- else -}}
    {{ default "default" .Values.metrics.serviceMonitor.name }}
{{- end -}}
{{- end -}}
