{{/*
Return the proper image name
{{ include "common.image" ( dict "context" . "image" .Values.path.to.the.image) }}
*/}}
{{- define "common.image" -}}
{{- $registryName := .image.registry -}}
{{- $repositoryName := .image.repository -}}
{{- $separator := ":" -}}
{{- $tag := .image.tag | toString -}}
{{- if and (empty $tag) (hasKey . "defaultTag") }}
    {{- $tag = .defaultTag | toString -}}
{{- end }}
{{- $termination := $tag -}}
{{- if and (not $registryName) (.context.Values.global).imageRegistry }}
    {{- $registryName = (.context.Values.global).imageRegistry -}}
{{- end -}}
{{- if (.context.Values.global).imageShortNames }}
    {{- $repositoryName = splitList "/" $repositoryName | last -}}
{{- end -}}
{{- if .image.digest }}
    {{- if $tag }}
        {{- $termination = printf "%s@%s" $tag .image.digest -}}
    {{- else }}
        {{- $separator = "@" -}}
        {{- $termination = .image.digest | toString -}}
    {{- end }}
{{- end -}}
{{- if $registryName }}
    {{- printf "%s/%s%s%s" $registryName $repositoryName $separator $termination -}}
{{- else -}}
    {{- printf "%s%s%s"  $repositoryName $separator $termination -}}
{{- end -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "common.imagePullSecrets" -}}
{{- $secrets := list -}}
{{- range concat (default (list) (.Values.global).imagePullSecrets) (default (list) .Values.imagePullSecrets) -}}
    {{- if kindIs "map" . }}
    {{- $secrets = append $secrets (dict "name" .name) }}
    {{- else }}
    {{- $secrets = append $secrets (dict "name" .) }}
    {{- end }}
{{- end -}}
{{- with $secrets -}}
imagePullSecrets: {{- toYaml . | nindent 2 }}
{{- end -}}
{{- end -}}
