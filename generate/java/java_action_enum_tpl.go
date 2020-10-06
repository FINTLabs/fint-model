package java

const ACTION_ENUM_TEMPLATE = 
`package {{ .Package }};

import java.util.Arrays;
import java.util.List;

public enum {{ .Name }} {
	{{ $c := sub (len .Classes) 1 -}}
	{{ range $i, $class := .Classes }}
	GET_{{ $class }},
	GET_ALL_{{ $class }},
	UPDATE_{{ $class }}{{ if ne $i $c }},{{ end }}
	{{- end }}
	;


    /**
     * Gets a list of all enums as string
     *
     * @return A string list of all enums
     */
    public static List<String> getActions() {
        return Arrays.asList(
                Arrays.stream({{ .Name }}.class.getEnumConstants()).map(Enum::name).toArray(String[]::new)
        );
    }

}

`
