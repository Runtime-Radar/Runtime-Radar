{{- define "postgresql.validateValues" -}}
{{- $messages := list -}}
{{- $messages = append $messages (include "common.tls.validateValues" .) -}}
{{- $messages = without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{- printf "\nVALUES VALIDATION:\n%s" $message | fail -}}
{{- end -}}
{{- end -}}
