{{/*
Set the ingress hostname
*/}}
{{- define "podinfo-ingress-hostname" -}}
{{- if .Values.ingress.host -}}
{{- print .Values.ingress.host -}}
{{- else -}}
{{- printf "%s.example.com" .Release.Name }}
{{- end -}}
{{- end }}
