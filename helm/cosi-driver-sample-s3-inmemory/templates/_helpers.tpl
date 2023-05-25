{{/*
Expand the name of the chart.
*/}}
{{- define "driver.name" }}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "driver.fullname" }}
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
Create chart name and version as used by the chart label.
*/}}
{{- define "driver.chart" }}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "driver.selectorLabels" }}
app.kubernetes.io/name: {{ include "driver.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "driver.labels" }}
helm.sh/chart: {{ include "driver.chart" . }}
{{ include "driver.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "driver.serviceAccountName" }}
{{- if .Values.serviceAccount.create }}
{{- default (include "driver.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the full image (repository and tag) for driver
*/}}
{{- define "driver.image" }}
{{- printf "%s/%s:%s" .Values.driver.image.registry .Values.driver.image.repository .Values.driver.image.tag }}
{{- end }}

{{/*
Create the full image (repository and tag) for sidecar
*/}}
{{- define "sidecar.image" }}
{{- printf "%s/%s:%s" .Values.sidecar.image.registry .Values.sidecar.image.repository .Values.sidecar.image.tag }}
{{- end }}

{{/*
Check if verbosity for driver is an integer, otherwise set to 5
*/}}
{{- define "driver.verbosity" }}
  {{- if (kindIs "int" .Values.driver.verbosity) }}
    {{- .Values.driver.verbosity }}
  {{- else }}
    {{- 5 }}
  {{- end }}
{{- end }}

{{/*
Check if verbosity for sidecar is an integer, otherwise set to 5
*/}}
{{- define "sidecar.verbosity" }}
  {{- if (kindIs "int" .Values.sidecar.verbosity) }}
    {{- .Values.sidecar.verbosity }}
  {{- else }}
    {{- 5 }}
  {{- end }}
{{- end }}
