{{/*
Return the proper RabbitMQ image name
*/}}
{{- define "rabbitmq.image" -}}
{{ include "common.image" (dict "context" . "image" .Values.image) }}
{{- end -}}

{{/*
Get the password key to be retrieved from RabbitMQ secret.
*/}}
{{- define "rabbitmq.secretPasswordKey" -}}
    {{- if and .Values.auth.existingSecret .Values.auth.existingSecretPasswordKey -}}
    {{- printf "%s" (tpl .Values.auth.existingSecretPasswordKey $) -}}
    {{- else -}}
        {{- printf "rabbitmq-password" -}}
    {{- end -}}
{{- end -}}

{{/*
Return the proper RabbitMQ plugin list
*/}}
{{- define "rabbitmq.plugins" -}}
{{- $plugins := .Values.plugins -}}
{{- if (eq (include "common.metrics.enabled" .) "true") -}}
{{- $plugins = printf "%s %s" $plugins .Values.metrics.plugins -}}
{{- end -}}
{{- printf "%s" $plugins | replace " " ", " -}}
{{- end -}}
