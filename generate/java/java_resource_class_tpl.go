package java

const RESOURCE_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ resourcePkg .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonSetter;

import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

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

    @Getter
    private final Map<String, List<Link>> links = Collections.synchronizedMap(new LinkedHashMap<>());

{{- if .Relations }}
    {{ range $i, $rel := .Relations }}
    public void add{{ upperCaseFirst $rel }}(Link link) {
        addLink("{{$rel}}", link);
    }
    {{ end }}
{{ end -}}
{{ if .Resources }}
    {{ range $att,$typ := .Resources }}
    @Getter
    private {{$typ}} {{$att}};
    @JsonSetter
    public void set{{ upperCaseFirst $att }}({{ $typ }} {{ $att }}) {
        this.{{$att}} = {{$att}};
    }
    @JsonIgnore
    @Override
    public void set{{ upperCaseFirst $att }}({{ baseType $typ }} {{$att}}) {
        this.{{$att}} = new {{$typ}}({{$att}});
    }
    {{ end }}
{{- end }}
}
`
