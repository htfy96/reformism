# Reformism
Utilities to empower Go's `{text/html}/template`.

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


{{ $x := slice 1 2 3 }}

{{ range $y := $x }}
    {{$y}},
{{end}}
// Result 1,2,3,

{{ $r0 := rng 5 }}
{{ range $e := $r0 }}
    {{$e}},
{{ end }}
// result: 0,1,2,3,4

{{ $r1 := rng 1,5 }}
{{ range $e := $r1 }}
    {{$e}},
{{ end }}
// result: 1,2,3,4,

{{ $r1_app := $r1 | append 5 6 }}
{{ $r1_app }}
// result 1,2,3,4,5,6

{{ $r2 := rng 10, 1, -3 }}
{{ range $e := $r2 }}
    {{$e}},
{{end}}
// result: 10,7,4,

{{ $el := "1,2,3" | split "," }}
{{$el}};
{{end}}
// result 1;2;3;

{{ $el := "1,2,3" | split "," | join ";" }}
{{ $el }}
{{end}}
// result 1;2;3

{{ $m := map "foo" 1 | map "bar" 2 }}
{{ range $k, $v := $m }}
    {{$k}}:{{$v}},
{{end}}
// Result: bar:1,foo:2,
		
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
 
#### slice
```
{{ slice 1 2 "abc" }}
```

make `[]interface{}`

#### rng
```
{{ rng {{count}} }} // 0, ..., count-1
{{ rng {{start}} {{end}} }}
{{ rng {{start}} {{end}} {{step}} }}
```
make `[]int` in given range

#### append
```
{{ append el1, el2, ..., slice }}
{{ slice | append el1, el2 }}
```

append elements to given slice

#### split
```
{{ split "separator" "str" }}
{{ "str" | split "separator" }}
```
split `string` by given separator to `[]string`

#### join
```
{{ "1,2,3" | split "," | join ";" }}
```

join `[]string` with given separator to `string`

#### map
```
{{ map "foo" 1 "bar" 2 }}
// Equivalent to 
{{ map "foo" 1 | map "bar" 2 }}
```

make `map[string]interface{}`

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