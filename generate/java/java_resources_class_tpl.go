package java

const RESOURCES_CLASS_TEMPLATE = `package {{ resourcePkg .Package }};

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.core.type.TypeReference;

import java.util.Collection;
import java.util.List;

import lombok.NoArgsConstructor;
import no.novari.fint.model.resource.AbstractCollectionResources;

@NoArgsConstructor
{{ if .Deprecated -}}
@Deprecated
{{ end -}}
public class {{ .Name }}Resources extends AbstractCollectionResources<{{.Name}}Resource> {

    public {{.Name}}Resources(Collection<{{.Name}}Resource> input) {
        super(input);
    }

    @JsonIgnore
    @Deprecated
    @Override
    public TypeReference<List<{{.Name}}Resource>> getTypeReference() {
        return new TypeReference<List<{{.Name}}Resource>>() {};
    }
}
`
