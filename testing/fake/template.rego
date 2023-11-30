package fake

default result = {}
default status = 200
default matched = false

result = {{.Result}} {{if .When}}{
    {{.When}}
}{{end}}

status = {{.Status}} {{if .When}}{
    {{.When}}
}{{end}}

{{if .When}}
matched {
    {{.When}}
}
{{end}}
