package java

const RESOURCE_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ resourcePkg .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonSetter;

import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import lombok.ToString;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import no.fint.model.{{ javaType .Stereotype }};
import no.fint.model.resource.FintLinks;
import no.fint.model.resource.Link;

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
public {{- if .Abstract }} abstract {{- end }} class {{ .Name }}Resource {{ if .Extends -}} extends {{ .Extends }}{{ if .ExtendsResource }}Resource{{ end }} {{ end -}} implements {{ javaType .Stereotype }}, FintLinks {

{{- if .Attributes }}
    // Attributes
    {{- if .Resources }}
    @JsonIgnore
    @Override
    public List<FintLinks> getNestedResources() {
        List<FintLinks> result = {{ if not .ExtendsResource }}FintLinks.{{end}}super.getNestedResources();
        {{- range $att := .Resources }}
        if ({{$att.Name}} != null) {
            result.add{{if $att.List}}All{{end}}({{$att.Name}});
        }
        {{- end }}
        return result;
    }
    {{- end }}
    {{- range $att := .Attributes }}
    {{- if not $att.Optional }}
    @NonNull
    {{- end }}
    private {{ javaType $att.Type | resource $.Resources | listFilt $att.List }} {{ $att.Name }};
    {{- end }}

{{- end }}

    // Relations
    @Getter
    private final Map<String, List<Link>> links = createLinks();

    {{- if .Relations }}
        {{ range $i, $rel := .Relations }}
    public void add{{ upperCaseFirst $rel }}(Link link) {
        addLink("{{$rel}}", link);
    }
        {{- end }}
    {{- end }}
}
`
