{{- define "cs.initSecretKey" -}}
{{- $value := "" -}}
{{- if empty .value -}}
    {{- if (pluck "lookup" (.context.Values.keys) ((.context.Values.global).keys) (dict "lookup" true) | first) -}}
    {{- $value = include "common.secrets.lookup" (dict "context" .context "secret" "keys" "key" .key) -}}
    {{- end -}}
{{- end -}}
{{- default (eq .value "INIT-DO-NOT-USE" | ternary (randBytes .size | b64dec | printf "%x") .value) $value -}}
{{- end -}}
