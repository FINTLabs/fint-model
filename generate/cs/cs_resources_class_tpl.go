package cs

const RESOURCES_CLASS_TEMPLATE = `using FINT.Model.Resource;

namespace {{.Namespace}}
{
    public class {{.Name}}Resources : AbstractCollectionResources<{{.Name}}Resource>
    {
    }
}
`
