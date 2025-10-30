{{- define "common.cs.container.app" -}}
- name: {{ include "common.name" . }}
  image: {{ include "common.cs.image" (dict "context" . "image" .Values.image) }}
  imagePullPolicy: {{ eq (include "common.cs.devMode.enabled" .) "true" | ternary "Always" (default "IfNotPresent" (.Values.image).pullPolicy) }}
  {{- with .Values.command }}
  command:
    {{- tpl (toYaml .) $ | nindent 4 }}
  {{- end }}
  env:
    {{- include "common.cs.container.app.env" . | nindent 4 }}
  envFrom:
    - configMapRef:
        name: {{ include "common.cs.configmapName" . }}
    {{- if (.Values.postgresql).enabled }}
    - secretRef:
        name: postgresql
    {{- end }}
    {{- if (.Values.redis).enabled }}
    - secretRef:
        name: redis
    {{- end }}
    {{- if (.Values.rabbitmq).enabled }}
    - secretRef:
        name: rabbitmq
    {{- end }}
    {{- if (.Values.clickhouse).enabled }}
    - secretRef:
        name: clickhouse
    {{- end }}
    {{- with .Values.envFrom }}
    {{- tpl (toYaml .) $ | nindent 4 }}
    {{- end }}
  {{- with .Values.containerPorts }}
  ports:
  {{- range $key, $val := . }}
  - name: {{ $key }}
    containerPort: {{ $val }}
  {{- end }}
  {{- end }}
  {{- with .Values.resources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- if (.Values.startupProbe).enabled }}
  startupProbe: {{- omit .Values.startupProbe "enabled" | toYaml | nindent 4 }}
  {{- end }}
  {{- if (.Values.livenessProbe).enabled }}
  livenessProbe: {{- omit .Values.livenessProbe "enabled" | toYaml | nindent 4 }}
  {{- end }}
  {{- if (.Values.readinessProbe).enabled }}
  readinessProbe: {{- omit .Values.readinessProbe "enabled" | toYaml | nindent 4 }}
  {{- end }}
  {{- if ne (include "common.cs.devMode.enabled" .) "true" }}
  {{- if (.Values.containerSecurityContext).enabled }}
  securityContext:
    {{- omit .Values.containerSecurityContext "enabled" | toYaml | nindent 4 }}
  {{- end }}
  {{- end }}
  {{- include "common.cs.volumeMounts" (dict "context" .) | nindent 2 }}
{{- end -}}

{{- define "common.cs.container.app.env" -}}
- name: LOG_LEVEL
  value: {{ include "common.cs.logLevel" . | quote }}
{{- if eq (include "common.cs.auth.enabled" .) "true" }}
- name: TOKEN_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "common.cs.keysSecretName" . }}
      key: token
{{- end }}
{{- if eq (include "common.cs.encryption.enabled" .) "true" }}
- name: ENCRYPTION_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "common.cs.keysSecretName" . }}
      key: encryption
{{- end }}
- name: GOPS_CONFIG_DIR
  value: "/tmp"
{{- if (.Values.rabbitmq).enabled }}
{{- with (.Values.rabbitmq).queue }}
- name: RABBIT_QUEUE
  value: {{ . | quote }}
{{- end }}
{{- with (.Values.rabbitmq).runtimeEventsQueue }}
- name: RABBIT_RUNTIME_EVENTS_QUEUE
  value: {{. | quote }}
{{- end }}
{{- with (.Values.rabbitmq).historyEventsQueue }}
- name: RABBIT_HISTORY_EVENTS_QUEUE
  value: {{. | quote }}
{{- end }}
{{- end }}
{{- with .Values.env }}
{{ tpl (toYaml .) $ }}
{{- end }}
{{- end -}}
