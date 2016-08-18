/*
Package reformism provides several utility functions for native text/template
*/
package reformism

import (
	"fmt"
	"reflect"
	"text/template"
)

// Pack represents packed arguments and original dot
type Pack struct {
	Origin interface{}
	Args   map[string]interface{}
}

// Witharg is used in pipe to pack argument with dot
func Witharg(k string, v interface{}, i interface{}) Pack {
	packT := reflect.TypeOf((*Pack)(nil)).Elem()
	if reflect.TypeOf(i) == packT {
		old := i.(Pack)
		old.Args[k] = v
		return old
	}
	return Pack{
		Origin: i,
		Args: map[string]interface{}{
			k: v,
		},
	}

}

// Done eats all pack passed to it and returns nil
func Done(Pack) interface{} {
	return nil
}

// Args extracts .Args field
func Args(p Pack) map[string]interface{} {
	return p.Args
}

// ArgCheckError may be raised in RequireArg
type ArgCheckError struct {
	detail string
}

// NewArgCheckError returns a new ArgCheckError instance from detailed message
func NewArgCheckError(s string) *ArgCheckError {
	return &ArgCheckError{
		detail: s,
	}
}

func (a ArgCheckError) Error() string {
	return a.detail
}

// RequireArg accepts packed dot(Pack), checks its validity, then returns the dot
func RequireArg(k string, trailingArgs ...interface{}) (interface{}, error) {
	if len(trailingArgs) != 1 && len(trailingArgs) != 2 {
		return nil, NewArgCheckError(`Invalid format. requireArg parameterName ["typeName"]`)
	}
	v := trailingArgs[len(trailingArgs)-1]

	if v, ok := v.(Pack); ok { // check whether last arg is Pack
		if _, ok := v.Args[k]; !ok { // check whether Pack contains arguments with name K
			return nil, NewArgCheckError(fmt.Sprintf("Required argument not found. Expected: %s, actual args: %v",
				k,
				v.Args))
		}
		if len(trailingArgs) == 2 { // check type
			if expectedTypeName, ok := trailingArgs[0].(string); ok {
				if reflect.TypeOf(v.Args[k]).Name() != expectedTypeName {
					return nil, NewArgCheckError(fmt.Sprintf("Unmatched type: Expected: %s, actual: %s",
						expectedTypeName,
						reflect.TypeOf(v.Args[k]).Name()))
				}
			} else {
				return nil, NewArgCheckError(fmt.Sprintf("The second argument of requireArg must be string! %v found",
					trailingArgs[0]))
			}
		}
		return trailingArgs[len(trailingArgs)-1], nil
	}
	return nil, NewArgCheckError("requireArg didn't receive argument modified by withArg")
}

func MakeSlice(args ...interface{}) []interface{} {
	return args
}

// MapError may be raised in MakeMap
type MapError struct {
	detail string
}

// NewMapError returns a new ArgCheckError instance from detailed message
func NewMapError(s string) *MapError {
	return &MapError{
		detail: s,
	}
}

func (a MapError) Error() string {
	return a.detail
}

func MakeMap(args ...interface{}) (map[string]interface{}, error) {
	if len(args) < 2 {
		return nil, NewMapError("arg num not required")
	}
	rawMap := make(map[string]interface{})
	if oldMap, ok := args[len(args)-1].(map[string]interface{}); ok {
		rawMap = oldMap
		args = args[:len(args)-1]
	}

	if len(args)%2 != 0 {
		return nil, NewMapError("arg should like key1 value1 key2 value2 ...")
	}
	for i := 0; i < len(args); i += 2 {
		if key, ok := args[i].(string); !ok {
			return nil, NewMapError("key should be string")
		} else {
			rawMap[key] = args[i+1]
		}

	}
	return rawMap, nil
}

// FuncsText is a FuncMap which can be passed as argument of .Func of text/template
var FuncsText = template.FuncMap{
	"arg":     Witharg,
	"require": RequireArg,
	"done":    Done,
	"args":    Args,
	"slice":   MakeSlice,
	"map":     MakeMap,
}

// FuncsHTML is a FuncMap which can be passed as argument of .Func of html/template
var FuncsHTML = FuncsText
