{{/*
Return the proper PostgreSQL image name
*/}}
{{- define "postgresql.image" -}}
{{ include "common.image" (dict "context" . "image" .Values.image) }}
{{- end -}}

{{/*
Return the path to the cert file.
*/}}
{{- define "postgresql.tlsCert" -}}
{{- printf "/etc/postgresql/certs/%s" (default "tls.crt" .Values.tls.certFilename) -}}
{{- end -}}

{{/*
Return the path to the cert key file.
*/}}
{{- define "postgresql.tlsCertKey" -}}
{{- printf "/etc/postgresql/certs/%s" (default "tls.key" .Values.tls.certKeyFilename) -}}
{{- end -}}

{{/*
Return the path to the CA cert file.
*/}}
{{- define "postgresql.tlsCACert" -}}
{{- printf "/etc/postgresql/certs/%s" (default "ca.crt" .Values.tls.certCAFilename) -}}
{{- end -}}

{{/*
Return the path to the CRL file.
*/}}
{{- define "postgresql.tlsCRL" -}}
{{- if .Values.tls.crlFilename -}}
{{- printf "/etc/postgresql/certs/%s" .Values.tls.crlFilename -}}
{{- end -}}
{{- end -}}

{{- define "postgresql.metrics.configmapName" -}}
{{- printf "%s-metrics" (include "common.fullname" .) -}}
{{- end -}}

{{/*
Get the admin-password key.
*/}}
{{- define "postgresql.adminPasswordKey" -}}
{{- if .Values.auth.existingSecret -}}
    {{- with .Values.auth.existingSecretAdminPasswordKey -}}
        {{- printf "%s" (tpl . $) -}}
    {{- end -}}
{{- else -}}
    {{- "postgres-password" -}}
{{- end -}}
{{- end -}}

{{/*
Get the password key.
*/}}
{{- define "postgresql.userPasswordKey" -}}
{{- $username := .Values.auth.username }}
{{- if or (empty $username) (eq $username "postgres") -}}
    {{- printf "%s" (include "postgresql.adminPasswordKey" .) -}}
{{- else -}}
    {{- if .Values.auth.existingSecret -}}
        {{- with .Values.auth.existingSecretPasswordKey -}}
            {{- printf "%s" (tpl . $) -}}
        {{- end -}}
    {{- else -}}
        {{- "password" -}}
    {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Get the probe command
*/}}
{{- define "postgresql.probeCommand" -}}
{{- $tlsEnabled := eq (include "common.tls.enabled" .) "true" -}}
{{- $customUser := .Values.auth.username -}}
{{- $database := .Values.auth.database -}}
- |
  exec pg_isready -U {{ default "postgres" $customUser | quote }} {{- if or $database $tlsEnabled }} -d "{{ if $database }}dbname={{ $database }}{{ end }} {{- if $tlsEnabled }} sslcert={{ include "postgresql.tlsCert" $ }} sslkey={{ include "postgresql.tlsCertKey" $ }}{{- end }}"{{- end }} -h 127.0.0.1 -p {{ $.Values.containerPorts.postgresql }}
{{- end -}}
