package cs

const RESOURCES_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

using FINT.Model.Resource;

namespace {{.Namespace}}
{
    public class {{.Name}}Resources extends AbstractCollectionResources<{{.Name}}Resource>
    {
    }
}
`
