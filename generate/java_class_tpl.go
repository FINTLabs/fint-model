package generate


const JAVA_CLASS_TEMPLATE = `package {{ .Package }};

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;

{{- if .Imports }}
{{ range $i := .Imports }}
import {{ $i }};
{{- end -}}
{{ end }}

@Data
@AllArgsConstructor
@NoArgsConstructor
@EqualsAndHashCode
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }} {{ if .Extends -}} extend {{ .Extends }} {{ end -}} {{ if .Identifiable -}} implements Identifiable {{ end -}} {

{{- if .Relations }}
	{{ $c := sub (len .Relations) 1 -}}
	public enum Relasjonsnavn {
		{{- range $i, $rel := .Relations }}
			{{ $rel }}{{if ne $i $c }},{{ end -}}
		{{ end }}
	}
{{ end -}}
{{ if .Attributes }}
	{{ range $att := .Attributes -}}
		private {{ javaType $att.Type}} {{ $att.Name }};
	{{ end -}}
{{ end -}}
{{ if .Identifiable }}
	@JsonIgnore
	@Override
	public String getId() {
		return this.get$FIXME$().getIdentifikatorverdi();
	}
{{ end -}}
}

`