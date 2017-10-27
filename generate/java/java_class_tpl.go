package java

const CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ .Package }};

import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.ToString;
import java.util.List;
import no.fint.model.*;

{{- if .Imports }}
{{ range $i := .Imports }}
import {{ $i }};
{{- end -}}
{{ end }}

@Data
{{ if .Extends -}}
@EqualsAndHashCode(callSuper=true)
@ToString(callSuper=true)
{{ else -}}
@EqualsAndHashCode
@ToString
{{ end -}}
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }} {{ if .Extends -}} extends {{ .Extends }} {{ end -}} implements {{ javaType .Stereotype }} {

{{- if .Relations }}
	{{ $c := sub (len .Relations) 1 -}}
	public enum Relasjonsnavn {
		{{- range $i, $rel := .Relations }}
			{{ $rel }}{{ if ne $i $c }},{{ end -}}
		{{ end }}
	}
{{ end -}}
{{ if .Attributes }}
	{{ range $att := .Attributes -}}
		private {{ javaType $att.Type | listFilt $att.List }} {{ $att.Name }};
	{{ end -}}
{{ end -}}
}
`
