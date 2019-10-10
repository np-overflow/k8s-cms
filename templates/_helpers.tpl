{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "k8s-cms.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "k8s-cms.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" $name .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "k8s-cms.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "k8s-cms.labels" -}}
{{- toYaml .Values.metadata.labels -}}
version: {{ .Chart.Version }}
instance: {{ .Release.Name }}
environment: {{ .Values.environment }}
app.kubernetes.io/name: {{ include "k8s-cms.name" . }}
app.kubernetes.io/fullname: {{ include "k8s-cms.fullname" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Label Selectors
Used by services and deployments to figure out which pods to manage/direct to.
*/}}
{{- define "k8s-cms.selectors" -}}
app.kubernetes.io/name: {{ include "k8s-cms.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Common annotations
Checksum annotations to force pods to reload configuration on update.
or secrets reload.
*/}}
{{- define "k8s-cms.annotations" -}}
{{- if.Values.metadata.annotations }}
{{- toYaml .Values.metadata.annotations -}}
{{- end -}}
checksum/{{ include "k8s-cms.fullname" . }}-secrets: {{ include (print $.Template.BasePath "/secrets.yml") . | sha256sum }}
checksum/{{ include "k8s-cms.fullname" . }}-config-env: {{ include (print $.Template.BasePath "/config/config-env.yml") . | sha256sum }}
checksum/{{ include "k8s-cms.fullname" . }}-config-file: {{ include (print $.Template.BasePath "/config/config-file.yml") . | sha256sum }}
{{- end -}}

