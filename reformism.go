/*
Reformism provides several utility functions for native text/template
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
	} else {
		return Pack{
			Origin: i,
			Args: map[string]interface{}{
				k: v,
			},
		}
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

type ArgCheckError struct {
	detail string
}

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
	} else {
		return nil, NewArgCheckError("requireArg didn't receive argument modified by withArg")
	}
}

// FuncsText is a FuncMap which can be passed as argument of .Func
var FuncsText = template.FuncMap{
	"arg":     Witharg,
	"require": RequireArg,
	"done":    Done,
	"args":    Args,
}

var FuncsHTML = FuncsText
