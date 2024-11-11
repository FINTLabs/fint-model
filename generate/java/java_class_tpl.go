package java

const CLASS_TEMPLATE = `package {{ .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.EqualsAndHashCode;
import lombok.ToString;
import lombok.Getter;
import java.util.Arrays;
import java.util.Collections;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import javax.validation.Valid;
import javax.validation.constraints.*;
import no.fint.model.FintMultiplicity;
import no.fint.model.felles.kompleksedatatyper.Identifikator;
import no.fint.model.{{ modelRename (javaType .Stereotype) }};
import no.fint.model.FintIdentifikator;
import no.fint.model.FintRelation;
{{- if .Imports -}}
{{ range $i := .Imports }}
import {{ $i }};
{{- end -}}
{{ end }}

import static no.fint.model.FintMultiplicity.ONE_TO_ONE;
import static no.fint.model.FintMultiplicity.ONE_TO_MANY;
import static no.fint.model.FintMultiplicity.NONE_TO_ONE;
import static no.fint.model.FintMultiplicity.NONE_TO_MANY;

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
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }} {{ if .Extends -}} extends {{ .Extends }} {{ end }} implements {{ modelRename (javaType .Stereotype) }}{{ if eq .Name "Identifikator"}}, FintIdentifikator{{ end }} {

{{- if .Relations }}
    {{ $c := sub (len .Relations) 1 -}}
    @Getter
    public enum Relasjonsnavn implements FintRelation {
        {{- range $i, $rel := .Relations }}
        {{ upperCase $rel.Name }}("{{ $rel.Name }}", "{{ $rel.Package }}.{{ $rel.Target }}", {{ resolveMultiplicity $rel.Multiplicity }}){{ if ne $i $c }},{{ else }};{{ end -}}
        {{ end }}
    
        private final String name;
        private final String packageName;
        private final FintMultiplicity multiplicity;

        private Relasjonsnavn(String name, String packageName, FintMultiplicity multiplicity) {
            this.name = name;
            this.packageName = packageName;
            this.multiplicity = multiplicity;
        }
    }
{{ end -}}
    
{{- if .Identifiable }}
    @JsonIgnore
    @Override
    private Map<String, FintIdentifikator> getIdentifikators() {
        Map<String, FintIdentifikator> identifikators = new HashMap<>();

        {{- if .ExtendsIdentifiable}}
        identifikators.putAll(super.getIdentifikators());
        {{- end}}

        {{- range $att := .Attributes }}
        {{- if eq $att.Type "Identifikator" }}
        identifikators.put("{{ $att.Name }}", this.{{ $att.Name }});
        {{- end }}
        {{- end }}

        return Collections.unmodifiableMap(identifikators);
    }
{{- end }}

{{- if .Relations }}
    @JsonIgnore
    private List<FintRelation> createRelations() {
        List<FintRelation> relations = new ArrayList<>();

        {{- if .ExtendsRelations }}
        relations.addAll(super.getRelations());
        {{- end }}

        relations.addAll(Arrays.asList(Relasjonsnavn.values()));

        return Collections.unmodifiableList(relations);
    }
{{- end }}

    public boolean isWriteable() {
        return this.writeable;
    }

    @JsonIgnore
    private final boolean writeable = {{ .Writable }};

{{- if .Relations }}
    @JsonIgnore
    private final List<FintRelation> relations = createRelations();
{{- end }}
{{- if .Attributes }}
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
