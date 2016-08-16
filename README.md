# Reformism
Add missing multiple arguments support for Go's `{text/html}/template`

## Usage
### Example
```go
import (
    "github.com/htfy96/reformism"
    "text/template"
    "os"
)

const template_text = `
{{define "foo"}}
	{{if $args := . | require "arg1" | require "arg2" "int" | args }}
	    {{with .Origin }} // Original dot
			{{.Bar}}
			{{$args.arg1}}
		{{ end }}
	{{ end }}
{{ end }}

{{ template "foo" . | arg "arg1" "Arg1" | arg "arg2" 42 }}
{{ template "foo" . | arg "arg1" "Arg1" | arg "arg2" "42" }} // will raise an error`


func main() {
    renderContext := map[string]string {
        "Bar": "bar",
    }
    t := template.Must(
        template.New("test_template").Funcs(
            reformism.FuncsText, // Use .FuncsHTML for html/template
        ).Parse(template_text))
    t.Execute(os.Stdout, renderContext)
}
```

### Docs
This package provides several utility functions for `{text/html}/template`, 
mappings to which are defined in `.FuncsText`(for `text/template`) and 
 `.FuncsHTML`(for `html/template`)

#### arg
```
{{ .Anything | arg "ArgName" ArgValue | ... }}
```

Convert anything to `Pack` type with argument stored.

#### require
```
// In template
{{ . | require "ArgName" | require "ArgName2" "typeName" | ... }}
```
Check `Pack` type's arguments, then returns `Pack` without any modification.

#### args
```
// In template
{{ if $args := . | ... | args }}
```

Extract `Args` Field of `Pack`. 

#### done
```
// In template 
{{ . | require "argname" | ... | done }}
```
Eats all data and returns nil.

#### Pack
```
type Pack struct {
	Origin interface{}
	Args   map[string]interface{}
}
```

## License
Apache. See `LICENSE` for more info.