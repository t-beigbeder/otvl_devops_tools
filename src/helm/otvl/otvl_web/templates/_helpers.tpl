{{/*
print a Pod environment from env value as dict
*/}}
{{- define "podenv-from-env" -}}
{{- if .Values.env -}}
{{- range $key, $value := .Values.env }}
- name: {{ printf "%s" $key }}
  value: {{ printf "%s" $value }}
{{- end }}
{{- end -}}
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
