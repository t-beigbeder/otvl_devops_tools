{{/*
Set the ingress hostname
*/}}
{{- define "netinfo-ingress-hostname" -}}
{{- if .Values.ingress.host -}}
{{- print .Values.ingress.host -}}
{{- else -}}
{{- printf "%s.example.com" .Release.Name }}
{{- end -}}
{{- end }}
