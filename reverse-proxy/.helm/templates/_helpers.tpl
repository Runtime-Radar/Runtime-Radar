{{- define "reverse-proxy.hosts" -}}
  {{- $hosts := append .Values.ingress.tls.hosts .Values.ingress.hostname }}
  {{- $hostname := include "common.cs.ownCsUrl" . | urlParse | pluck "hostname" | first }}
  {{- append $hosts $hostname | compact | uniq | dict "hosts" | toYaml }}
{{- end -}}

{{- define "reverse-proxy.tls.secretName" -}}
{{- default (printf "%s-ingress-crt" (include "common.cs.basename" .)) ((.Values.ingress).tls).existingSecret -}}
{{- end -}}
