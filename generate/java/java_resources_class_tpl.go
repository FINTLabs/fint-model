package java

const RESOURCES_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

package {{ resourcePkg .Package }};

import com.fasterxml.jackson.core.type.TypeReference;

import java.util.List;

import no.fint.model.resource.AbstractCollectionResources;

public class {{ .Name }}Resources extends AbstractCollectionResources<{{.Name}}Resource> {

    @Override
    public TypeReference<List<{{.Name}}Resource>> getTypeReference() {
        return new TypeReference<List<{{.Name}}Resource>>() {};
    }
}
`
