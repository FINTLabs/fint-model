package cs

const RESOURCE_CLASS_TEMPLATE = `// Built from tag {{ .GitTag }}

using System;
using System.Collections.Generic;
using Newtonsoft.Json;
using FINT.Model.Resource;

{{- if .Using }}
{{ range $u := .Using }}
using {{ $u }};
{{- end -}}
{{ end }}

namespace {{ .Namespace }}
{

    public {{- if .Abstract }} abstract {{- end }} class {{ .Name }}Resource {{ if .Extends -}} : {{ .Extends }}{{ if .ExtendsResource }}Resource{{ end }} {{ end }}
    {

    {{ if .Attributes }}
        {{ range $att := .Attributes -}}
        public {{ csType $att.Type $att.Optional | resource $.Resources | listFilt $att.List }} {{ upperCaseFirst $att.Name }} { get; set; }
        {{ end -}}
    {{ end }}

    {{- if not .ExtendsResource }}
        {{if .Abstract}}protected{{else}}public{{end}} {{.Name}}Resource()
        {
            Links = new Dictionary<string, List<Link>>();
        }

        [JsonProperty(PropertyName = "_links")]
        public Dictionary<string, List<Link>> Links { get; private set; }

        protected void AddLink(string key, Link link)
        {
            if (!Links.ContainsKey(key))
            {
                Links.Add(key, new List<Link>());
            }
            Links[key].Add(link);
        }
     {{ end -}}

        {{- if .Relations }}
            {{ range $i, $rel := .Relations }}

        public void Add{{ upperCaseFirst $rel.Name }}(Link link)
        {
            AddLink("{{$rel.Name}}", link);
        }
            {{- end }}
        {{- end }}
    }
}
`
