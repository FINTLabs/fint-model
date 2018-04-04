package java

const RESOURCE_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ resourcePkg .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonSetter;

import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.util.ArrayList;
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
    public static {{.Name}}Resource create({{.Name}} other) {
        if (other == null) {
            return null;
        }
        if (other instanceof {{.Name}}Resource) {
            return other;
        }
        {{.Name}}Resource result = new {{.Name}}Resource();
        {{- range $att := .AllAttributes }}
        result.set{{ upperCaseFirst $att.Name }}(other.get{{ upperCaseFirst $att.Name }}());
        {{- end }}
        return result;
    }

{{- if .Resources }}
    // Resources
    @JsonIgnore
    @Override
    public List<FintLinks> getNestedResources() {
        List<FintLinks> result = new ArrayList<>();
        {{- range $att,$typ := .Resources }}
        if ({{ getter $att }} != null) {
            result.add{{ listAdder $typ}}({{ getter $att | assignResource $typ }});
        }
        {{- end }}
        return result;
    }
    {{ range $att,$typ := .Resources }}
    @JsonSetter
    @Override
    public void set{{ upperCaseFirst $att }}({{ baseType $typ }} {{$att}}) {
        super.set{{ upperCaseFirst $att }}({{ assignResource $typ $att }});
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
