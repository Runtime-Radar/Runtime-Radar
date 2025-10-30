{{/*
Merge a list of values.
Merge precedence is consistent with http://masterminds.github.io/sprig/dicts.html#merge-mustmerge
Usage:
{{ include "common.utils.merge" ( dict "values" (list .Values.path.to.the.Value1 .Values.path.to.the.Value2) "context" $ ) }}
*/}}
{{- define "common.utils.merge" -}}
{{- $dst := dict -}}
{{- range .values -}}
{{- $dst = merge $dst . -}}
{{- end -}}
{{ $dst | toYaml }}
{{- end -}}
