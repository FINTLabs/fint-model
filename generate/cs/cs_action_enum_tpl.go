package cs

const ACTION_ENUM_TEMPLATE = `// Built from tag {{ .GitTag }}

using System;

namespace {{ .Namespace }}
{
	public enum {{ .Name }}
    {
	{{ $c := sub (len .Classes) 1 -}}
	{{ range $i, $class := .Classes }}
	GET_{{ $class }},
	GET_ALL_{{ $class }},
	UPDATE_{{ $class }}{{ if ne $i $c }},{{ end }}
	{{- end }}
    }
}
`
