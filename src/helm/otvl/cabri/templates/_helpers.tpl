{{/*
print a Pod environment from var
*/}}
{{- define "podenv-from-var" -}}
{{- if .Values.env -}}
{{- range $key, $value := .Values.env }}
- name: {{ printf "%s" $key }}
  value: {{ printf "%s" $value }}
{{- end }}
{{- end -}}
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
