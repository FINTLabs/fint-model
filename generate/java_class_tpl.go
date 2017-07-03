package generate

const JAVA_CLASS_TEMPLATE = `package {{ .Package }};

{{ if .Identifiable }}
import com.fasterxml.jackson.annotation.JsonIgnore;
import no.fint.model.relation.Identifiable;
{{ end }}
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
{{ if .Extends -}}
@EqualsAndHashCode(callSuper=false)
{{ else -}}
@AllArgsConstructor
@NoArgsConstructor
@EqualsAndHashCode
{{ end -}}
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }} {{ if .Extends -}} extends {{ .Extends }} {{ end -}} {{ if .Identifiable -}} implements Identifiable {{ end -}} {

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
		return this.getSystemId().getIdentifikatorverdi();
	}
{{ end -}}
}

`
