{{/*
print a Pod environment from cabri_webapi values
*/}}
{{- define "podenv-from-cabri-webapi" -}}
- name: CABRI_SYNC_WEB_USER
  value: {{ printf "%s" .Values.cabri_webapi.user }}
- name: CABRI_SYNC_WEB_PASSWORD
  value: {{ printf "%s" .Values.cabri_webapi.password }}
- name: CABRI_SYNC_DSS
  value: {{ printf "%s" .Values.cabri_webapi.dss }}
- name: CABRI_SYNC_CONFIG
  value: {{ printf "%s" .Values.cabri_webapi.config }}
- name: CABRI_SYNC_CONTENT
  value: {{ printf "%s" .Values.cabri_webapi.content }}
- name: CABRI_SYNC_ASSETS
  value: {{ printf "%s" .Values.cabri_webapi.assets }}
{{- end -}}

{{/*
print a list of DSS from dsss value as list of strings
*/}}
{{- define "dsss-args-from-dsss" -}}
{{- range $value := .Values.dsss }}
- {{ printf "%s" $value }}
{{- end }}
{{- end -}}

{{/*
print a list of auth_users from apr_credentials list value as list of strings
*/}}
{{- define "auth-users-from-credentials" -}}
{{- range $value := .Values.apr_credentials }}
- {{ printf "%s" $value }}
{{- end }}
{{- end -}}

{{/*
Set the ingress hostname
*/}}
{{- define "ingress-hostname" -}}
{{- if .Values.ingress.host -}}
{{- print .Values.ingress.host -}}
{{- else -}}
{{- printf "%s.example.com" .Release.Name }}
{{- end -}}
{{- end }}
