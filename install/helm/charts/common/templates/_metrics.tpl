{{/*
Return the proper metrics image name
*/}}
{{- define "common.metrics.image" -}}
{{ include "common.image" (dict "context" . "image" (.Values.metrics).image) }}
{{- end -}}

{{/*
Metrics enabled and external host is empty
*/}}
{{- define "common.metrics.enabled" -}}
{{- if and (.Values.metrics).enabled (not (.Values.metrics).externalHost) -}}
    {{- true -}}
{{- end -}}
{{- end -}}
