package reformism

import (
	"bytes"
	"text/template"
	"strings"
	"testing"
)

type testCase struct {
	template       string
	argument       interface{}
	expectedResult string
	hasError       bool
}

var testCases = []testCase{
	testCase{
		template: `
		{{define "foo"}}
		{{if $args := . | require "arg1" | require "arg2" "int" | args }}
		{{with .Origin }}
			{{.Bar}}
			{{$args.arg1}}
		{{ end }}
		{{ end }}
		{{ end }}

		{{template "foo" . | arg "arg1" "test" | arg "arg2" 42}}
		`,
		argument: map[string]string{
			"Bar": "bar",
		},
		expectedResult: "bartest",
		hasError:       false,
	},
	testCase{
		template: `
		{{define "foo"}}
		{{if $args := . | require "arg1" | require "arg2" "string" | args }}
		{{with .Origin }}
			{{.Bar}}
			{{$args.arg1}}
		{{ end }}
		{{ end }}
		{{ end }}

		{{template "foo" . | arg "arg1" "test" | arg "arg2" 42}}
		`,
		argument: map[string]string{
			"Bar": "bar",
		},
		expectedResult: "bartest",
		hasError:       true,
	},
}

func removeWhite(s string) string {
	toRemove := []string {
		"\n", " ", "\t",
	}
	for _, r := range toRemove {
		s = strings.Replace(s, r, "", -1)
	}
	return s
}

func runTestCase(t *testing.T, tc testCase) {
	temp := template.Must(
		template.New("test_template").Funcs(
			FuncsText,
	).Parse(tc.template))
	buf := new(bytes.Buffer)
	err := temp.Execute(buf, tc.argument)
	if (err != nil) != tc.hasError {
		t.Errorf("haserror status unexpected. Expected %v, actual error %v", tc.hasError, err)
	}
	if !tc.hasError && removeWhite(buf.String()) != removeWhite(tc.expectedResult) {
		t.Errorf("Unexpected result. Expected: %s, Actual: %s", tc.expectedResult, buf.String())
	}
}

func TestAll(t *testing.T) {
	for _, tc := range testCases {
		runTestCase(t, tc)
	}
}
