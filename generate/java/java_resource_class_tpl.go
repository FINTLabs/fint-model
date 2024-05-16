package java

const RESOURCE_CLASS_TEMPLATE = `package {{ resourcePkg .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;

import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import javax.validation.Valid;
import javax.validation.constraints.*;

import no.fint.model.felles.kompleksedatatyper.Identifikator;
import no.fint.model.resource.FintLinks;
import no.fint.model.{{ resourcePackageRename (javaType .Stereotype) }};
import no.fint.model.resource.Link;
import no.fint.model.FintIdentifikator;

{{- if .Imports -}}
{{ range $i := .Imports }}
import {{ resource $.Resources $i | extends $.ExtendsResource $.Extends }};
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
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }}Resource {{ if .Extends -}} extends {{ .Extends }}{{ if .ExtendsResource }}Resource{{ end }} {{ end -}} implements {{ implementInterfaces (resourceRename (javaType .Stereotype)) }} {

{{- if .Attributes }}
    // Attributes
    {{- if .Resources }}
    @JsonIgnore
    @Override
    public List<FintLinks> getNestedResources() {
        List<FintLinks> result = {{ if not .ExtendsResource }}{{ superResource .Stereotype }}.{{end}}super.getNestedResources();
        {{- range $att := .Resources }}
        if ({{$att.Name}} != null) {
            result.add{{if $att.List}}All{{end}}({{$att.Name}});
        }
        {{- end }}
        return result;
    }
    {{- end }}
    {{- range $att := .Attributes }}
    {{- if $att.Deprecated }}
    @Deprecated
    {{- end }}
    {{- if not $att.Optional }}
    {{ if $att.List }}@NotEmpty{{ else if eq "string" $att.Type }}@NotBlank{{ else }}@NotNull{{ end }}
    {{- end }}
    private {{ javaType $att.Type | resource $.Resources | validFilt $att.Type | listFilt $att.List }} {{ $att.Name }};
    {{- end }}

{{- end }}

{{- if .Identifiable }}
    @JsonIgnore
    public Map<String, FintIdentifikator> getIdentifikators() {
        Map<String, FintIdentifikator> identifikators = new HashMap<>();

    {{- if .ExtendsIdentifiable}}
        identifikators.putAll(super.getIdentifikators());
    {{- end}}

    {{- range $att := .Attributes}}
    {{- if eq $att.Type "Identifikator"}}
        identifikators.put("{{ $att.Name }}", this.{{ $att.Name }});
    {{- end}}
    {{- end}}
    
        return identifikators;
    }
{{- end }}

    // Relations
    @Getter
    private final Map<String, List<Link>> links = createLinks();

    {{- if .Relations }}
        {{ range $i, $rel := .Relations }}

    {{- if $rel.Deprecated }}
    @Deprecated
    {{- end }}
    @JsonIgnore
    public List<Link> get{{ upperCaseFirst $rel.Name }}() {
        return getLinks().getOrDefault("{{$rel.Name}}", Collections.emptyList()); 
    }
    {{- if $rel.Deprecated }}
    @Deprecated
    {{- end }}
    public void add{{ upperCaseFirst $rel.Name }}(Link link) {
        addLink("{{$rel.Name}}", link);
    }
        {{- end }}
    {{- end }}
}
`
