package generate

const CS_CLASS_TEMPLATE = `
{{ if .Using }}
{{ range $u := .Using }}
using {{ $u }};
{{- end -}}
{{ end }}

namespace {{ .Namespace }}
{
	public class {{ .Name }}
	{
		{{- if .Relations }}
		{{ $c := sub (len .Relations) 1 -}}
		public enum Relasjonsnavn
        {
		{{- range $i, $rel := .Relations }}
			{{ $rel }}{{if ne $i $c }},{{ end -}}
		{{ end }}
        }
        {{ end }}
	{{ if .Attributes }}
		{{ range $att := .Attributes -}}
			public {{ $att.Type}} {{ $att.Name }} { get; set; }
		{{ end -}}
	{{ end }}
	}
}
`
