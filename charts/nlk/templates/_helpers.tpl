{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "nlk.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "nlk.fullname" -}}
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
Create a default fully qualified nlk name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "nlk.nlk.fullname" -}}
{{- printf "%s-%s" (include "nlk.fullname" .) .Values.nlk.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified nlk service name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "nlk.nlk.service.name" -}}
{{- default (include "nlk.nlk.fullname" .) .Values.serviceNameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "nlk.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "nlk.labels" -}}
helm.sh/chart: {{ include "nlk.chart" . }}
{{ include "nlk.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "nlk.selectorLabels" -}}
{{- if .Values.nlk.selectorLabels -}}
{{ toYaml .Values.nlk.selectorLabels }}
{{- else -}}
app.kubernetes.io/name: {{ include "nlk.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
{{- end -}}

{{/*
Expand the name of the configmap.
*/}}
{{- define "nlk.configName" -}}
{{- if .Values.nlk.customConfigMap -}}
{{ .Values.nlk.customConfigMap }}
{{- else -}}
{{- default (include "nlk.fullname" .) .Values.nlk.config.name -}}
{{- end -}}
{{- end -}}

{{/*
Expand service account name.
*/}}
{{- define "nlk.serviceAccountName" -}}
{{- default (include "nlk.fullname" .) .Values.nlk.serviceAccount.name -}}
{{- end -}}

{{- define "nlk.tag" -}}
{{- default .Chart.AppVersion .Values.nlk.image.tag -}}
{{- end -}}

{{/*
Expand image name.
*/}}
{{- define "nlk.image" -}}
{{- if .Values.nlk.image.digest -}}
{{- printf "%s/%s@%s" .Values.nlk.image.registry .Values.nlk.image.repository .Values.nlk.image.digest -}}
{{- else -}}
{{- printf "%s/%s:%s" .Values.nlk.image.registry .Values.nlk.image.repository (include "nlk.tag" .) -}}
{{- end -}}
{{- end -}}

{{- define "nlk.prometheus.serviceName" -}}
{{- printf "%s-%s" (include "nlk.fullname" .) "prometheus-service"  -}}
{{- end -}}
