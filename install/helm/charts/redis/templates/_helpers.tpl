{{/*
Return the proper Redis image name
*/}}
{{- define "redis.image" -}}
{{ include "common.image" (dict "context" . "image" .Values.image) }}
{{- end -}}

{{/*
Get the password key to be retrieved from Redis&reg; secret.
*/}}
{{- define "redis.secretPasswordKey" -}}
{{- if and .Values.auth.existingSecret .Values.auth.existingSecretPasswordKey -}}
{{- printf "%s" (tpl .Values.auth.existingSecretPasswordKey $) -}}
{{- else -}}
{{- printf "redis-password" -}}
{{- end -}}
{{- end -}}

{{/*
Return the path to the cert file.
*/}}
{{- define "redis.tlsCert" -}}
{{- printf "/etc/redis/certs/%s" (default "tls.crt" .Values.tls.certFilename) -}}
{{- end -}}

{{/*
Return the path to the cert key file.
*/}}
{{- define "redis.tlsCertKey" -}}
{{- printf "/etc/redis/certs/%s" (default "tls.key" .Values.tls.certKeyFilename) -}}
{{- end -}}

{{/*
Return the path to the CA cert file.
*/}}
{{- define "redis.tlsCACert" -}}
{{- printf "/etc/redis/certs/%s" (default "ca.crt" .Values.tls.certCAFilename) -}}
{{- end -}}
