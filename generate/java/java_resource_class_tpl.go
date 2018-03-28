package java

const RESOURCE_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ resourcePkg .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonSetter;

import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import no.fint.model.resource.FintLinks;
import no.fint.model.resource.Link;

import {{ .Package }}.{{ .Name }};
{{- if .Imports -}}
{{ range $i := .Imports }}
import {{ $i }};
{{- end -}}
{{ end }}

@NoArgsConstructor
@EqualsAndHashCode(callSuper=true)
@ToString(callSuper=true)
public class {{ .Name }}Resource extends {{ .Name }} implements FintLinks {
{{- if .Resources }}
    // Resources
    {{- range $att,$typ := .Resources }}
    @Getter
    private {{$typ}} {{$att}};
    {{- end }}

    {{- range $att,$typ := .Resources }}
    @JsonSetter
    public void set{{ upperCaseFirst $att }}({{ $typ }} _{{ $att }}) {
        this.{{$att}} = _{{$att}};
    }
    @JsonIgnore
    @Override
    public void set{{ upperCaseFirst $att }}({{ baseType $typ }} _{{$att}}) {
        this.{{$att}} = {{ assignResource $typ $att }};
    }
    {{- end }}
{{- end }}

    // Links
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
