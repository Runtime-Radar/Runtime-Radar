{{- define "common.cs.basename" -}}
{{- printf "cs" -}}
{{- end -}}

{{- define "common.cs.configmapName" -}}
{{- printf "%s-config" (include "common.cs.basename" .) -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "common.cs.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "common.basename" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return own cs url
*/}}
{{- define "common.cs.ownCsUrl" -}}
{{- $url := default (.Values.global).ownCsUrl .Values.ownCsUrl -}}
{{- required "global.ownCsUrl argument is required for installation" $url -}}
{{- end -}}

{{/*
Return central cs url
*/}}
{{- define "common.cs.centralCsUrl" -}}
{{- $url := default (.Values.global).centralCsUrl .Values.centralCsUrl -}}
{{- if eq (include "common.cs.isChildCluster" .) "true" -}}
{{- required "global.centralCsUrl argument is required for installation" $url -}}
{{- else -}}
{{- default (include "common.cs.ownCsUrl" .) $url -}}
{{- end -}}
{{- end -}}

{{/*
Return cs version
*/}}
{{- define "common.cs.csVersion" -}}
{{- default (.Values.global).csVersion .Values.csVersion | default "v0.0.0" -}}
{{- end -}}

{{/*
Return keys secret name
*/}}
{{- define "common.cs.keysSecretName" -}}
{{- default "cs-keys" ((.Values.global).keys).existingSecret -}}
{{- end -}}

{{- define "common.cs.grafana.address" -}}
{{- default (.Values.grafana).externalHost ((.Values.global).grafana).externalHost | default (printf "%s://grafana.%s.svc.cluster.local:3000" (include "common.cs.http-scheme" .) .Release.Namespace) | trimSuffix "/" }}
{{- end -}}
