package java

const CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ .Package }};

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.EqualsAndHashCode;
import lombok.ToString;
import lombok.NonNull;
import java.util.List;
import no.fint.model.{{ javaType .Stereotype }};
{{- if .Imports -}}
{{ range $i := .Imports }}
import {{ $i }};
{{- end -}}
{{ end }}

@Data
@NoArgsConstructor
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
            {{ upperCase $rel }}{{ if ne $i $c }},{{ end -}}
        {{ end }}
    }
{{ end -}}
{{ if .Attributes }}
    {{- range $att := .Attributes }}
    {{- if not $att.Optional }}
    @NonNull
    {{- end }}
    private {{ javaType $att.Type | listFilt $att.List }} {{ $att.Name }};
    {{- end }}
{{- end }}
}
`
