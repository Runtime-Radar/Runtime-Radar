{{/*
Return the proper image name
{{ include "common.cs.image" ( dict "context" . "image" .Values.path.to.the.image) }}
*/}}
{{- define "common.cs.image" -}}
{{- with (.context.Values.global).imageRegistry -}}
{{- $_ := set $.image "registry" . -}}
{{- end -}}
{{- $tag := coalesce (.context.Values.global).imageTag (.context.Values.global).csVersion .context.Chart.Version -}}
{{- include "common.image" (dict "defaultTag" $tag | mergeOverwrite .) -}}
{{- end -}}
