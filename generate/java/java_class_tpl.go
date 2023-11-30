package java

const CLASS_TEMPLATE = `package {{ .Package }};

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.EqualsAndHashCode;
import lombok.ToString;
import java.util.List;
import javax.validation.Valid;
import javax.validation.constraints.*;
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
{{ if .Deprecated -}}
@Deprecated
{{ end -}}
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }} {{ if .Extends -}} extends {{ .Extends }} {{ end -}} implements {{ javaType .Stereotype }} {

{{- if .Relations }}
    {{ $c := sub (len .Relations) 1 -}}
    public enum Relasjonsnavn {
        {{- range $i, $rel := .Relations }}
            {{ upperCase $rel.Name }}("{{ $rel.Target }}", "{{ $rel.Multiplicity }}"){{ if ne $i $c }},{{ end -}}
        {{ end }}
	
		private final String typeName;
        private final String multiplicity;

        private Relasjonsnavn(String typeName, String multiplicity) {
            this.typeName = typeName;
            this.multiplicity = multiplicity;
        }

        public String getTypeName() {
            return typeName;
        }

        public String getMultiplicity() {
            return multiplicity;
        }
    }
{{ end -}}
{{ if .Attributes }}
    {{- range $att := .Attributes }}
    {{- if $att.Deprecated }}
    @Deprecated
    {{- end }}
    {{- if not $att.Optional }}
    {{ if $att.List }}@NotEmpty{{ else if eq "string" $att.Type }}@NotBlank{{ else }}@NotNull{{ end }}
    {{- end }}
    private {{ javaType $att.Type | validFilt $att.Type | listFilt $att.List }} {{ $att.Name }};
    {{- end }}
{{- end }}
}
`
